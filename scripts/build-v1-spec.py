"""Derive openapi-v1.json (the V1-only, Speakeasy-annotated spec) from original-openapi.json.

Keeps every /v1/* path so the generated SDK covers the whole V1 API, prunes the
V0 surface, and annotates the subset that maps onto Terraform resources and data
sources.

Not every endpoint has a declarative shape. Imperative actions (POST .../sync,
POST .../tasks/{id}:complete and friends) and multipart uploads stay in the SDK
but are deliberately left unannotated -- Terraform has nowhere to put them.

Run `make spec` (dump + build + overlay) rather than calling this directly.
"""

import argparse
import copy
import json
import sys
from collections import defaultdict

ALPHA_DESCRIPTION = (
    "ALPHA: the 1.x line targets the Lumos V1 API and may introduce breaking "
    'changes without prior notice. Pin an exact version. For production, stay '
    'on 0.10.x (version = "~> 0.10").'
)

# entity -> {operation: (path, method)}. Entity names become lumos_<snake_case>.
RESOURCES = {
    "App": {
        "create": ("/v1/apps", "post"),
        "read": ("/v1/apps/{app_uuid}", "get"),
        "update": ("/v1/apps/{app_uuid}", "put"),
        "delete": ("/v1/apps/{app_uuid}", "delete"),
    },
    "Entitlement": {
        "create": ("/v1/entitlements", "post"),
        "read": ("/v1/entitlements/{entitlement_uuid}", "get"),
        "update": ("/v1/entitlements/{entitlement_uuid}", "patch"),
        "delete": ("/v1/entitlements/{entitlement_uuid}", "delete"),
    },
    "Workflow": {
        "create": ("/v1/workflows", "post"),
        "read": ("/v1/workflows/{workflow_id}", "get"),
        "update": ("/v1/workflows/{workflow_id}", "put"),
        "delete": ("/v1/workflows/{workflow_id}", "delete"),
    },
}

DATA_SOURCES = {
    "Apps": ("/v1/apps", "get"),
    "AppStoreCatalog": ("/v1/apps/appstore_catalog", "get"),
    "AppEntitlements": ("/v1/apps/{app_uuid}/entitlements", "get"),
    "AppEntitlementAppStoreCatalog": ("/v1/apps/{app_uuid}/entitlements/appstore_catalog", "get"),
    "AppStoreSettings": ("/v1/apps/{app_uuid}/appstore-settings", "get"),
    "AccountEntitlements": ("/v1/accounts/{account_uuid}/entitlements", "get"),
    "CustomAttributeDefinitions": ("/v1/custom-attribute-definitions", "get"),
    "EntitlementAccounts": ("/v1/entitlements/{entitlement_uuid}/accounts", "get"),
    "EntitlementAppStoreSettings": ("/v1/entitlements/{entitlement_uuid}/appstore_settings", "get"),
    "Integrations": ("/v1/integrations", "get"),
    "Integration": ("/v1/integrations/{app_class_id}", "get"),
    "IntegrationConnectionSchema": ("/v1/integrations/{app_class_id}/connection-schema", "get"),
    "KnowledgeContexts": ("/v1/knowledge/context", "get"),
    "KnowledgeContext": ("/v1/knowledge/context/{knowledge_context_id}", "get"),
    "Tasks": ("/v1/tasks", "get"),
    "Task": ("/v1/tasks/{task_id}", "get"),
    "TaskActions": ("/v1/tasks/{task_id}/actions", "get"),
    "Users": ("/v1/users", "get"),
    "CurrentUser": ("/v1/users/current", "get"),
    "User": ("/v1/users/{user_id}", "get"),
    "UserAccounts": ("/v1/users/{user_id}/accounts", "get"),
    "UserRoles": ("/v1/users/{user_id}/roles", "get"),
    "Workflows": ("/v1/workflows", "get"),
    "WorkflowSteps": ("/v1/workflow-steps", "get"),
}

# Singleton workflow bindings: GET reads the bound workflow, PUT binds one.
# There is no unbind endpoint, so these resources have no delete.
#
# All six share the Workflow and WorkflowAssociationInput schemas, and
# x-speakeasy-entity binds at the schema level, so each association needs its own
# cloned pair before it can be its own Terraform resource.
#
# entity -> (path, owning path parameter)
WORKFLOW_ASSOCIATIONS = {
    "AppApprovalWorkflow": ("/v1/apps/{app_uuid}/approval-workflow", "app_uuid"),
    "AppProvisioningWorkflow": (
        "/v1/apps/{app_uuid}/default-provisioning-workflow",
        "app_uuid",
    ),
    "AppDeprovisioningWorkflow": (
        "/v1/apps/{app_uuid}/default-deprovisioning-workflow",
        "app_uuid",
    ),
    "EntitlementApprovalWorkflowOverride": (
        "/v1/entitlements/{entitlement_uuid}/approval-workflow-override",
        "entitlement_uuid",
    ),
    "EntitlementProvisioningWorkflow": (
        "/v1/entitlements/{entitlement_uuid}/provisioning-workflow",
        "entitlement_uuid",
    ),
    "EntitlementDeprovisioningWorkflow": (
        "/v1/entitlements/{entitlement_uuid}/deprovisioning-workflow",
        "entitlement_uuid",
    ),
}


# Path parameter that identifies the resource, surfaced to Terraform as `id`.
ID_PARAM = {
    "App": "app_uuid",
    "Entitlement": "entitlement_uuid",
    "Workflow": "workflow_id",
}

# Request-body fields that cannot change in place; changing them replaces the resource.
#
# app_class_id decides which integration the app is, so changing it is a different
# app, not an edit. It is the successor to the old top-level app_id.
#
# The old app_instance_identifier force-new is deliberately absent: the field moved
# inside the free-form `settings` object and is no longer a typed property, so there
# is nothing to hang a plan modifier on.
FORCE_NEW = {
    "ConnectIntegrationInput": ["app_class_id"],
}


def op_of(spec, path, method):
    try:
        return spec["paths"][path][method]
    except KeyError:
        sys.exit(f"missing operation: {method.upper()} {path}")


def body_schemas(op):
    out = []
    for content in (op.get("requestBody") or {}).get("content", {}).values():
        ref = (content.get("schema") or {}).get("$ref")
        if ref:
            out.append(ref.split("/")[-1])
    return out


def response_schema(op):
    for code in ("200", "201", "202"):
        resp = op.get("responses", {}).get(code)
        if not resp:
            continue
        ref = (resp.get("content", {}).get("application/json", {}).get("schema") or {}).get("$ref")
        return ref.split("/")[-1] if ref else None
    return None


def set_response_ref(op, name):
    for code in ("200", "201", "202"):
        resp = op.get("responses", {}).get(code)
        if resp and "$ref" in (
            resp.get("content", {}).get("application/json", {}).get("schema") or {}
        ):
            resp["content"]["application/json"]["schema"]["$ref"] = f"#/components/schemas/{name}"
            return


def sanitize_descriptions(node, counter):
    """Break up Go-template delimiters that appear inside descriptions.

    Speakeasy embeds descriptions into its own .stmpl template source before
    parsing it, so a literal ``}}`` in a description closes a template action
    that was never opened and generation dies with

        failed to parse template: provider/data_source.go.stmpl:43:
        unexpected "{" in command

    Inserting a space is meaning-preserving for the JSON snippets these
    descriptions carry, and keeps the text readable in the generated docs.
    """
    if isinstance(node, dict):
        for key, value in node.items():
            if key == "description" and isinstance(value, str):
                if "}}" in value or "{{" in value:
                    node[key] = value.replace("}}", "} }").replace("{{", "{ {")
                    counter.append(1)
            else:
                sanitize_descriptions(value, counter)
    elif isinstance(node, list):
        for value in node:
            sanitize_descriptions(value, counter)


def set_request_ref(op, name):
    content = (op.get("requestBody") or {}).get("content", {}).get("application/json")
    if content and "$ref" in (content.get("schema") or {}):
        content["schema"]["$ref"] = f"#/components/schemas/{name}"


def collect_refs(node, acc):
    if isinstance(node, dict):
        ref = node.get("$ref")
        if isinstance(ref, str) and ref.startswith("#/components/schemas/"):
            acc.add(ref.split("/")[-1])
        for value in node.values():
            collect_refs(value, acc)
    elif isinstance(node, list):
        for value in node:
            collect_refs(value, acc)


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--src", default="original-openapi.json")
    parser.add_argument("--out", default="openapi-v1.json")
    parser.add_argument(
        "--data-sources",
        default=None,
        help="comma-separated subset of data-source entities (debugging aid)",
    )
    args = parser.parse_args()

    global DATA_SOURCES
    if args.data_sources is not None:
        keep = {n for n in args.data_sources.split(",") if n}
        DATA_SOURCES = {k: v for k, v in DATA_SOURCES.items() if k in keep}

    src = json.load(open(args.src))
    spec = copy.deepcopy(src)

    # 1. V1 paths only.
    spec["paths"] = {p: v for p, v in src["paths"].items() if p.startswith("/v1/")}
    spec["info"]["description"] = ALPHA_DESCRIPTION
    spec.setdefault("servers", [{"url": "https://api.lumos.com"}])
    spec["x-speakeasy-retries"] = {
        "strategy": "backoff",
        "backoff": {
            "initialInterval": 500,
            "maxInterval": 30000,
            "maxElapsedTime": 60000,
            "exponent": 1.5,
        },
        "statusCodes": ["5XX"],
        "retryConnectionErrors": True,
    }

    schemas = spec.setdefault("components", {}).setdefault("schemas", {})

    # 2. Which schema does each entity bind to? Clone when two entities share one.
    claims = defaultdict(list)
    for entity, ops in RESOURCES.items():
        path, method = ops["read"]
        claims[response_schema(op_of(spec, path, method))].append(entity)
    for entity, (path, method) in DATA_SOURCES.items():
        claims[response_schema(op_of(spec, path, method))].append(entity)

    renames = {}
    for schema_name, entities in claims.items():
        if schema_name is None or len(entities) < 2:
            continue
        for entity in entities[1:]:
            clone = f"{entity}Output"
            if clone in schemas:
                clone = f"{entity}Result"
            schemas[clone] = copy.deepcopy(schemas[schema_name])
            renames[entity] = clone
            print(f"  cloned {schema_name} -> {clone} (shared by {entities})", file=sys.stderr)

    # 3. Annotate operations and bind entity schemas.
    bound = {}
    for entity, ops in RESOURCES.items():
        for action, (path, method) in ops.items():
            op = op_of(spec, path, method)
            op["x-speakeasy-entity-operation"] = f"{entity}#{action}"
            for schema_name in body_schemas(op):
                bound.setdefault(schema_name, entity)
            id_param = ID_PARAM.get(entity)
            for param in op.get("parameters", []):
                if id_param and param.get("name") == id_param and param.get("in") == "path":
                    param["x-speakeasy-name-override"] = "id"
        path, method = ops["read"]
        target = renames.get(entity) or response_schema(op_of(spec, path, method))
        if renames.get(entity):
            for _, (p, m) in ops.items():
                set_response_ref(op_of(spec, p, m), renames[entity])
        if target:
            bound.setdefault(target, entity)

    for entity, (path, method) in DATA_SOURCES.items():
        op = op_of(spec, path, method)
        op["x-speakeasy-entity-operation"] = f"{entity}#read"
        if entity in renames:
            set_response_ref(op, renames[entity])
            bound.setdefault(renames[entity], entity)
        else:
            target = response_schema(op)
            if target:
                bound.setdefault(target, entity)

    for schema_name, entity in bound.items():
        if schema_name in schemas:
            schemas[schema_name]["x-speakeasy-entity"] = entity

    for schema_name, fields in FORCE_NEW.items():
        props = schemas.get(schema_name, {}).get("properties", {})
        for field in fields:
            if field in props:
                props[field]["x-speakeasy-param-force-new"] = True

    # 3b. Workflow bindings: clone the shared schemas so each binding is its own entity.
    for entity, (path, owner) in WORKFLOW_ASSOCIATIONS.items():
        out_name, in_name = entity, f"{entity}Input"
        schemas[out_name] = copy.deepcopy(schemas["Workflow"])
        schemas[in_name] = copy.deepcopy(schemas["WorkflowAssociationInput"])
        schemas[out_name]["x-speakeasy-entity"] = entity
        schemas[in_name]["x-speakeasy-entity"] = entity

        read_op = op_of(spec, path, "get")
        write_op = op_of(spec, path, "put")
        set_response_ref(read_op, out_name)
        set_response_ref(write_op, out_name)
        set_request_ref(write_op, in_name)
        read_op["x-speakeasy-entity-operation"] = f"{entity}#read"
        # One idempotent PUT serves as both create and update.
        write_op["x-speakeasy-entity-operation"] = [f"{entity}#create", f"{entity}#update"]

        # Rebinding to a different app or entitlement is a different association.
        for op in (read_op, write_op):
            for param in op.get("parameters", []):
                if param.get("name") == owner and param.get("in") == "path":
                    param["x-speakeasy-param-force-new"] = True

    # 4. Keep Speakeasy's template parser out of the descriptions.
    touched = []
    sanitize_descriptions(spec, touched)

    # 5. Drop schemas the V1 surface no longer reaches.
    reachable = set()
    collect_refs(spec["paths"], reachable)
    frontier = set(reachable)
    while frontier:
        current = frontier.pop()
        nested = set()
        collect_refs(schemas.get(current, {}), nested)
        new = nested - reachable
        reachable |= new
        frontier |= new
    dropped = [name for name in schemas if name not in reachable]
    for name in dropped:
        del schemas[name]

    with open(args.out, "w") as handle:
        handle.write(json.dumps(spec, indent=2, sort_keys=True) + "\n")

    annotated = sum(
        1
        for item in spec["paths"].values()
        for op in item.values()
        if isinstance(op, dict) and "x-speakeasy-entity-operation" in op
    )
    total = sum(
        1
        for item in spec["paths"].values()
        for m, op in item.items()
        if m in ("get", "post", "put", "patch", "delete")
    )
    print(
        f"wrote {args.out}: {len(spec['paths'])} paths, {total} operations "
        f"({annotated} mapped to Terraform, {total - annotated} SDK-only), "
        f"{len(schemas)} schemas kept, {len(dropped)} pruned, "
        f"{len(touched)} descriptions sanitised",
        file=sys.stderr,
    )


if __name__ == "__main__":
    main()
