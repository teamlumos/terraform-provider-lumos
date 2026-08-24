# V1 Alpha Provider Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Publish a V1-API-only Terraform provider as `teamlumos/lumos` version `1.0.0-alpha.1`, generated from a long-lived `v1` branch, referencing zero V0 endpoints.

**Architecture:** A long-lived `v1` branch in this same repo carries its own `overlay-v1.yaml` (which strips every non-`/v1` path and adds the Speakeasy entity annotations the V1 endpoints currently lack) plus its own `.speakeasy/gen.yaml` version. The provider source address stays `registry.terraform.io/teamlumos/lumos` — unchanged — so the version number alone distinguishes V1 (1.x) from V0 (0.x). Terraform never auto-selects a prerelease, so existing `~> 0.10` customers are unaffected. The existing tag-triggered `terraform_publish.yaml` needs no changes because git tags are not branch-scoped.

**Tech Stack:** Speakeasy 1.700.2+ (`terraform` target, gen.yaml configVersion 2.0.0), Go 1.24.11, terraform-plugin-framework 1.17.0, tfplugindocs 0.24.0, GoReleaser v2, GitHub Actions.

**Spec:** No separate design doc — the design was agreed in conversation and is inlined under "Design Decisions" below. This plan is the spec of record.

**Shape:** Five tasks. Tasks 1-2 prepare the branch and every generation input; Task 3 runs a single generation and verifies everything it produced; Task 4 is the functional gate against the live API; Task 5 releases. `make provider` runs code generation *and* `go generate ./...` for docs in one pass (`Makefile:11-20`), so there is exactly one generation step in this plan and all of its inputs — overlay, version, docs template — are staged before it.

## Global Constraints

- Provider source address is `registry.terraform.io/teamlumos/lumos` and **must never change**. `main.go:32` already hardcodes it; do not edit that line.
- Provider version for this release is exactly `1.0.0-alpha.1`. Subsequent alphas increment the prerelease identifier (`1.0.0-alpha.2`, …).
- Zero V0 paths, resources, data sources, docs, or examples may survive on the `v1` branch. The only permitted API paths are the six under `/v1`.
- Go version floor: `go 1.24.11` (from `go.mod:3`). Do not bump it.
- GPG signing fingerprint: `F7FFD8D8D816AE3F8A6B6A0C9B9EA29C4ED23991` (from `.github/workflows/terraform_publish.yaml:31`). Unchanged.
- Alpha warning copy — use this exact wording everywhere it appears:
  `ALPHA: the 1.x line targets the Lumos V1 API and may introduce breaking changes without prior notice. Pin an exact version. For production, stay on 0.10.x (version = "~> 0.10").`
- Never run `make overlay` on the `v1` branch — it writes to `overlay.yaml` (the V0 overlay) and would clobber it. Task 2 adds a separate target for the V1 overlay.
- All work happens on branch `v1`. Never push V1 changes to `master`. The one exception is Task 2 Step 6, a single-flag change on `master`.
- **The V0 line must remain independently releasable.** Master keeps every V0 artifact and `terraform_publish.yaml` is unchanged, so `git tag v0.11.0` on master still ships a V0 release. Task 2 Step 6 makes this robust rather than merely true-by-default; Task 5 Step 8 proves it.

---

## Design Decisions

Recorded here so an executor does not have to re-derive them.

1. **Same source address, 1.x version.** `teamlumos/lumos` is at `0.10.4` after 59 published releases, so it never left 0.x and a breaking 1.0 rewrite is semver-legitimate. Reusing the address means EA adopters never rewrite a `required_providers` block at GA.
2. **V0 and V1 deliberately cannot coexist** in one configuration. This is a confirmed requirement, not a limitation to work around. Terraform loads one version per provider address; a customer needing both uses two root modules with separate state.
3. **Public prerelease, not a private repo.** Public visibility is acceptable given loud alpha labeling. This is also the cheaper path: the private-repo alternative needs a new repo, three copied workflows, two secrets, a manual artifact-distribution channel, and a GA migration back into this repo.
4. **Accepted cost:** the registry will display `1.0.0-alpha.1` as "latest" on the `teamlumos/lumos` page, and per the [registry FAQ](https://developer.hashicorp.com/terraform/registry/faq) this **cannot be prevented** — "The Terraform Registry's version sorting is strictly SemVer compliant." Mitigated by the alpha banner in `docs/index.md`, which the registry renders from the displayed tag. Existing `~> 0.10` customers are not affected.

---

## Deferred Scope: Knowledge Context

Not a task in this plan. Recorded so the omission is deliberate and reversible.

The six V1 paths do not all make good Terraform resources today:

| V1 path | Ops | Alpha.1 scope | Reason |
|---|---|---|---|
| `/v1/apps` | GET, POST | ✅ `lumos_app` + `lumos_apps` | Clean. `AppWithCustomAttributes.id` is a required string. |
| `/v1/apps/{app_uuid}` | GET, PUT, DELETE | ✅ `lumos_app` | Full CRUD, first V1 endpoint with a real DELETE. |
| `/v1/integrations` | GET | ✅ `lumos_integrations` | Data source only. |
| `/v1/integrations/{app_id}` | GET | ✅ `lumos_integration` | Data source only. |
| `/v1/knowledge/context` | GET, POST | ⚠️ **Deferred** | See below. |
| `/v1/knowledge/context/{id}` | GET, PATCH, DELETE | ⚠️ **Deferred** | See below. |

**Knowledge Context is deferred to `1.0.0-alpha.2`** for two concrete reasons found in the spec:

- **Multipart file upload cannot round-trip.** `Body_createKnowledgeContext` and `Body_updateKnowledgeContext` are `multipart/form-data` with `file` typed `{"type": "string", "format": "binary"}`. But `KnowledgeContextOutput` returns `file_url`, not `file`. There is no way for `Read` to recover the attribute `Create` sent, so every `terraform plan` would either show permanent drift or need `file` marked write-only — a modeling decision, not a generation detail.
- **`DELETE /v1/knowledge/context/{knowledge_context_id}` declares no success response.** Its `responses` map contains only `400, 401, 409, 422, 429, 500`. This is an upstream spec bug in `lumos/backend/python/services/apiserver/` and must be fixed there before the delete can generate correctly.

**To pick it up later,** resolve both blockers in `lumos/backend/python/services/apiserver/`, then:

1. Add the two paths to the `KEEP` list in Task 2 Step 1 and annotate `KnowledgeContext#create/read/update/delete`, tagging `KnowledgeContextOutput` and `ListKnowledgeContextOutput` with `x-speakeasy-entity`.
2. Re-run Tasks 2 and 3, then tag `v1.0.0-alpha.2`.

On the `file` attribute, options ranked:
- **Recommended:** drop `file` from the Terraform resource entirely. Model `lumos_knowledge_context` over `title`, `description`, `generic_tags`, `product_scopes` only, and leave file upload to the API or CLI. Terraform is a poor fit for binary payloads it cannot read back.
- Accept a `file_path` attribute marked `RequiresReplace` and excluded from `Read`, acknowledging that changing the file's *contents* without changing its path produces no diff.
- Add a server-side content hash to `KnowledgeContextOutput` so drift becomes detectable, then model `file` normally.

If you want Knowledge Context in alpha.1 after all, fold the above into Tasks 2-3 before starting, and budget time for the file-attribute modeling.

---

## File Structure

On branch `v1` only, except where noted.

**Created:**
- `overlay-v1.yaml` — the V1 overlay. Strips all non-`/v1` paths, adds `x-speakeasy-entity` / `x-speakeasy-entity-operation` annotations, rewrites `info.description`. Generated by `speakeasy overlay compare`, not hand-written.
- `openapi-v1.json` — the post-overlay V1 spec, committed for review/diffing. Mirrors how `openapi.json` is committed on master.
- `templates/index.md.tmpl` — tfplugindocs template carrying the alpha banner. Needed because `docs/` is **not** Speakeasy-tracked (0 `docs/` entries in `.speakeasy/gen.lock` `trackedFiles`, vs 27 for `examples/` and 80 for `internal/provider/`) and is produced by `go generate` via the directive at `main.go:18`. **Must exist before Task 3 generates**, or the first generation produces a bannerless index.

**Modified:**
- `.speakeasy/workflow.yaml:8` — overlay location `./overlay.yaml` → `./overlay-v1.yaml`, plus a separate registry namespace.
- `.speakeasy/gen.yaml:36` — `terraform.version: 0.10.4` → `1.0.0-alpha.1`.
- `.speakeasy/gen.lock` — clear the `persistentEdits` block (lines 10-13) and the `trackedFiles` map so the first V1 generation establishes a fresh baseline. See Task 1 for why.
- `.genignore` — remove the two stale V0 example entries; drop the `registration.go` entry.
- `internal/sdk/internal/hooks/registration.go` — revert to the empty generated stub.
- `Makefile` — `VERSION` override plus a `--match 'v1.*'` tag filter. **On `master`, the matching `--match 'v0.*'` filter** (Task 2 Step 6).
- `README.md` — alpha banner, V1 examples.
- `tests/client_test.go` — **hard compile break.** `tests/client_test.go:33` calls `client.AppStore.GetAppStoreApp` with `operations.GetAppStoreAppRequest`, and `:46` calls `client.Core.CreateApp` with `shared.AppInputCreate`. All four symbols disappear when V0 is stripped, and the `Core` and `AppStore` sub-clients cease to exist. Must be rewritten against V1 (Task 3 Step 7). `make provider` will not catch this — it only builds `./internal/provider/...`.
- `.goreleaser.yml` — one line, `prerelease: false` (Task 5). `v1` branch only; master's release behavior stays byte-identical to today.

**Deleted:**
- All 20 files under `docs/` (V0 resource/data-source/guide docs). Not auto-removed by generation.
- `examples/resources/lumos_access_policy/resource.tf` and `examples/resources/lumos_requestable_permission/resource.tf` — hand-written V0 examples, currently protected by `.genignore`.
- `internal/sdk/internal/hooks/requestable_permissions_hook.go` — a hand-written V0 hook that appends `include_inherited_configs=false` to URLs containing `requestable_permissions`. Note: it only imports `net/http` and `strings`, so it would still *compile* on V1 — it would just be dead code. It must go together with the `registration.go` reference to it.
- The remaining V0 files under `examples/` and `internal/provider/` — deleted by Speakeasy generation, since they are in `trackedFiles`. Task 3 verifies this rather than assuming it.

---

### Task 1: Prepare the `v1` branch

Cut the branch, remove every hand-maintained V0 artifact, and reset the generation baseline. Deliverable: a branch that still compiles and carries no V0 content that generation cannot remove on its own.

**Why the baseline reset matters.** `.speakeasy/gen.lock:10-13` contains:

```yaml
persistentEdits:
  generation_id: c748b837-001f-4465-aa56-531d05d0f07d
  pristine_commit_hash: 23d14bec32f4c6f6341df608eda8e48526d0c4f1
  pristine_tree_hash: 2833644eff3ec4aa11f7e0fc81047ae1764efe79
```

plus ~1700 lines of per-file `last_write_checksum` entries under `trackedFiles`. With `persistentEdits.enabled: "true"` in `gen.yaml:33-34`, Speakeasy diffs the working tree against that pristine baseline to find manual edits and re-applies them onto freshly generated output. That baseline is the **V0** pristine tree. On the `v1` branch, the first generation could read "51 V0 provider files and 27 V0 examples are gone" as manual edits worth preserving, and resurrect them.

The exact re-application semantics are inferred from the field names, so **Task 3 Step 3 verifies empirically instead of trusting the inference.**

**Files:**
- Create: branch `v1`
- Modify: `.speakeasy/gen.lock` (clear `persistentEdits` + `trackedFiles`), `.speakeasy/gen.yaml:33-34`, `.genignore`, `internal/sdk/internal/hooks/registration.go`
- Delete: `docs/` (20 files), 2 V0 example files, `internal/sdk/internal/hooks/requestable_permissions_hook.go`

**Interfaces:**
- Produces: a `v1` branch that compiles, with no hand-maintained V0 artifacts and a cleared generation baseline. Task 2 authors the generation inputs on top of it.

- [ ] **Step 1: Cut the branch from the current commit**

```bash
cd /Users/grahamseyffert/develop/terraform-provider-lumos
git checkout -b v1 1b72390
git status
```

Expected: `On branch v1`. The uncommitted `README.md` change and untracked `openapi.json` carry over — that is fine, they get handled in Tasks 2 and 3.

- [ ] **Step 2: Delete the hand-maintained V0 artifacts**

```bash
git rm -r docs/
git rm examples/resources/lumos_access_policy/resource.tf \
       examples/resources/lumos_requestable_permission/resource.tf
git rm internal/sdk/internal/hooks/requestable_permissions_hook.go
```

Expected: 23 files staged for deletion. `docs/` is gone entirely — Task 3 regenerates it.

- [ ] **Step 3: Revert `registration.go` to the empty generated stub**

It currently registers the V0 hook you just deleted, so the package will not compile until this is done. Write exactly:

```go
// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package hooks

/*
 * This file is only ever generated once on the first generation and then is free to be modified.
 * Any hooks you wish to add should be registered in the InitHooks function. Feel free to define them
 * in this file or in separate files in the hooks package.
 */

func initHooks(h *Hooks) {
	// Add hooks by calling h.register{SDKInit/BeforeRequest/AfterSuccess/AfterError}Hook
	// with an instance of a hook that implements that specific Hook interface
	// Hooks are registered per SDK instance, and are valid for the lifetime of the SDK instance
}
```

- [ ] **Step 4: Rewrite `.genignore`**

All three current entries are now obsolete: the two example files are deleted, and `registration.go` carries no manual edits worth protecting (its own header says Speakeasy only generates it once, so it does not need `.genignore` protection). Replace the file contents with:

```
# Intentionally empty on the v1 branch.
# The V0 entries (lumos_access_policy / lumos_requestable_permission examples and
# hooks/registration.go) were removed with their files. Add entries here only for
# files carrying hand-written content that Speakeasy would otherwise overwrite.
```

- [ ] **Step 5: Clear the persistent-edits baseline**

Disable persistent edits for the first generation only:

```bash
python3 - <<'EOF'
p = '.speakeasy/gen.yaml'
s = open(p).read()
s = s.replace('  persistentEdits:\n    enabled: "true"', '  persistentEdits:\n    enabled: "false"')
open(p, 'w').write(s)
EOF
grep -A1 "persistentEdits" .speakeasy/gen.yaml
```

Expected: `enabled: "false"`.

Then strip the stale baseline and tracked-file map out of the lock file:

```bash
python3 - <<'EOF'
p = '.speakeasy/gen.lock'
lines = open(p).read().split('\n')
out, skip = [], None
for ln in lines:
    if ln in ('persistentEdits:', 'trackedFiles:'):
        skip = ln
        continue
    if skip and (ln.startswith('  ') or ln == ''):
        continue
    skip = None
    out.append(ln)
open(p, 'w').write('\n'.join(out))
EOF
grep -c "persistentEdits\|trackedFiles" .speakeasy/gen.lock
```

Expected: `0`.

- [ ] **Step 6: Verify the branch still compiles**

This is the test for this task. The hooks package was just edited by hand, so a compile failure here means Step 3 was wrong.

```bash
go build ./... && go vet ./internal/sdk/internal/hooks/
```

Expected: both succeed with no output. If you see `undefined: RequestablePermissionsHook`, Step 3 did not take.

Note that `go build ./...` still passes at this point because the V0 SDK is still present — `tests/client_test.go` does not break until after Task 3 generates. That is expected and handled in Task 3 Step 7.

- [ ] **Step 7: Commit**

```bash
git add -A
git commit -m "v1: strip V0 artifacts and reset generation baseline

Removes hand-maintained V0 docs, examples, and the requestable-permissions
hook, and clears the V0 persistent-edits baseline so the first V1 generation
starts from a clean slate rather than treating removed V0 files as manual edits."
```

---

### Task 2: Author every generation input

The V1 spec, the overlay derived from it, the Speakeasy wiring, the version plumbing, and the docs banner template. Deliverable: everything Task 3's single `make provider` invocation needs, staged and asserted before anything is generated.

**Approach for the overlay:** do **not** hand-write ~40 JSONPath actions. This repo's established workflow (`Makefile:22-26`) is the reverse: edit the spec, then derive the overlay with `speakeasy overlay compare`. Follow it.

The annotation pattern to mirror is `overlay.yaml:28-45` (path-level `x-speakeasy-entity-operation` and `x-speakeasy-name-override`) and `overlay.yaml:280-325` (schema-level `x-speakeasy-entity`).

**Files:**
- Create: `overlay-v1.yaml`, `openapi-v1.json`, `templates/index.md.tmpl`
- Modify: `.speakeasy/workflow.yaml:8`, `Makefile` (on `v1`), `Makefile:6` (on `master`)

**Interfaces:**
- Consumes: the `v1` branch from Task 1.
- Produces: `overlay-v1.yaml` referenced by `.speakeasy/workflow.yaml`, yielding a spec whose only paths are the four in-scope V1 paths, with entities `App`, `Apps`, `Integration`, `Integrations`; plus a docs template Task 3's generation will render.

- [ ] **Step 1: Build the V1 spec from the original**

```bash
python3 - <<'PY'
import json, collections

spec = json.load(open('original-openapi.json'))

KEEP = ['/v1/apps', '/v1/apps/{app_uuid}', '/v1/integrations', '/v1/integrations/{app_id}']
spec['paths'] = collections.OrderedDict(
    (p, spec['paths'][p]) for p in KEEP if p in spec['paths']
)

spec['info']['description'] = (
    'ALPHA: the 1.x line targets the Lumos V1 API and may introduce breaking '
    'changes without prior notice. Pin an exact version. For production, stay '
    'on 0.10.x (version = "~> 0.10").'
)

# HTTP exponential backoff retries — same policy as the V0 overlay.
spec['x-speakeasy-retries'] = {
    'strategy': 'backoff',
    'backoff': {'initialInterval': 500, 'maxInterval': 30000,
                'maxElapsedTime': 60000, 'exponent': 1.5},
    'statusCodes': ['5XX'],
    'retryConnectionErrors': True,
}

P = spec['paths']

# --- App: full CRUD on /v1/apps + /v1/apps/{app_uuid} ---
P['/v1/apps']['post']['x-speakeasy-entity-operation'] = 'App#create'
P['/v1/apps']['get']['x-speakeasy-entity-operation'] = 'Apps#read'
P['/v1/apps/{app_uuid}']['get']['x-speakeasy-entity-operation'] = 'App#read'
P['/v1/apps/{app_uuid}']['put']['x-speakeasy-entity-operation'] = 'App#update'
P['/v1/apps/{app_uuid}']['delete']['x-speakeasy-entity-operation'] = 'App#delete'

# app_uuid is the connected-instance id (AppWithCustomAttributes.id), NOT the
# catalog id that ConnectIntegrationInput.app_id carries. Surface it as `id`.
for method in ('get', 'put', 'delete'):
    P['/v1/apps/{app_uuid}'][method]['parameters'][0]['x-speakeasy-name-override'] = 'id'

# --- Integration: read-only data sources ---
P['/v1/integrations']['get']['x-speakeasy-entity-operation'] = 'Integrations#read'
P['/v1/integrations/{app_id}']['get']['x-speakeasy-entity-operation'] = 'Integration#read'

S = spec['components']['schemas']
S['AppWithCustomAttributes']['x-speakeasy-entity'] = 'App'
S['ConnectIntegrationInput']['x-speakeasy-entity'] = 'App'
S['Page_AppWithCustomAttributes_']['x-speakeasy-entity'] = 'Apps'
S['IntegrationCatalogOutput']['x-speakeasy-entity'] = 'Integration'
S['Page_IntegrationCatalogOutput_']['x-speakeasy-entity'] = 'Integrations'

json.dump(spec, open('openapi-v1.json', 'w'), indent=2)
print('paths:', list(spec['paths']))
PY
```

Expected output: `paths: ['/v1/apps', '/v1/apps/{app_uuid}', '/v1/integrations', '/v1/integrations/{app_id}']`

- [ ] **Step 2: Assert the spec is V1-only and fully annotated**

This is the primary test for this task — run it before generating anything.

```bash
python3 - <<'PY'
import json, sys
spec = json.load(open('openapi-v1.json'))
paths = list(spec['paths'])
fail = []

bad = [p for p in paths if not p.startswith('/v1/')]
if bad:
    fail.append(f'non-V1 paths present: {bad}')
if len(paths) != 4:
    fail.append(f'expected 4 paths, got {len(paths)}: {paths}')

ops = {}
for p, item in spec['paths'].items():
    for m, op in item.items():
        if m in ('get', 'post', 'put', 'patch', 'delete'):
            ops[f'{m.upper()} {p}'] = op.get('x-speakeasy-entity-operation')
missing = [k for k, v in ops.items() if not v]
if missing:
    fail.append(f'operations missing x-speakeasy-entity-operation: {missing}')

if 'ALPHA' not in spec['info']['description']:
    fail.append('info.description is missing the ALPHA warning')

for k, v in sorted(ops.items()):
    print(f'  {k:45s} -> {v}')
print()
if fail:
    print('FAIL'); [print(' -', f) for f in fail]; sys.exit(1)
print('PASS: 4 V1 paths, 7 annotated operations, alpha warning present')
PY
```

Expected: seven annotated operations listed, then `PASS`. If it prints FAIL, fix Step 1 and re-run.

**Note on leftover schemas.** Step 1 prunes `paths` but leaves `components.schemas` intact, so every V0 schema (`AppStoreApp`, `AccessPolicyOutput`, `Page_User_`, …) is still in `openapi-v1.json`. This is intentional and harmless for correctness: Terraform resources come from `x-speakeasy-entity` annotations, which only the five V1 schemas carry, so no V0 resource can be generated. It does mean Speakeasy will emit unused `internal/sdk/models/shared/*.go` types — bloat in the binary, not a defect. If the generated tree looks unreasonably large after Task 3, prune unreferenced schemas in Step 1 as a follow-up; do not treat it as a blocker for alpha.1.

- [ ] **Step 3: Derive the overlay**

```bash
speakeasy overlay compare --schemas original-openapi.json --schemas openapi-v1.json > overlay-v1.yaml
head -20 overlay-v1.yaml
grep -c "x-speakeasy-entity-operation" overlay-v1.yaml
```

Expected: a valid overlay header, and `7` entity-operation annotations.

- [ ] **Step 4: Wire the Speakeasy workflow to the V1 overlay and its own namespace**

Edit `.speakeasy/workflow.yaml`. Two changes — the overlay location and the registry namespace:

```yaml
workflowVersion: 1.0.0
speakeasyVersion: latest
sources:
    FastAPI:
        inputs:
            - location: ./original-openapi.json
        overlays:
            - location: ./overlay-v1.yaml
        registry:
            location: registry.speakeasyapi.dev/lumos/lumos/fast-api-v1
targets:
    lumos-terraform-provider-target:
        target: terraform
        source: FastAPI
```

**Why the namespace change.** Both branches otherwise publish their compiled source spec to `registry.speakeasyapi.dev/lumos/lumos/fast-api`, and `.speakeasy/workflow.lock` records it tagged `latest`. Generation itself reads local files (`./original-openapi.json` plus the branch's overlay), so V0 codegen cannot break from this — but generating on `v1` would push a V1-filtered spec to that shared `latest` tag, and generating on master would push the V0-filtered one back. Anything else in your Speakeasy workspace consuming `fast-api:latest` would see it flip-flop. No other consumer exists in the local workspace and the `lumos` monorepo has no `.speakeasy` config, but the Speakeasy cloud workspace may hold language-SDK targets that are not visible from here. If the dashboard confirms nothing else consumes `fast-api`, this change is optional; it costs nothing either way.

- [ ] **Step 5: Add the `VERSION` override and v1 tag filter to the Makefile**

`Makefile:3-9` derives the version from `git describe --tags --abbrev=0`. Two problems on the `v1` branch: before any `v1.*` tag exists it returns `v0.10.4` and would stamp the provider `0.10.4`; and you want to generate and verify *before* tagging a release, which is circular.

It also needs a **tag-pattern filter**. `git describe --tags --abbrev=0` finds the most recent tag *reachable from HEAD*. A `v1.*` tag on this branch is unreachable from master today, so master is fine — but that safety rests on the branches never merging, which is a weak guarantee for something that fails silently. Pinning each branch to its own tag series makes the two lines immune to each other regardless of merge history.

Replace `Makefile:2-9` with:

```makefile
provider: overlay-v1.yaml
	@echo "Syncing gen.yaml version..."
	@VERSION="$(VERSION)"; \
	if [ -z "$$VERSION" ]; then \
	  VERSION=$$(git describe --tags --abbrev=0 --match 'v1.*' 2>/dev/null | sed 's/^v//'); \
	fi; \
	if [ -n "$$VERSION" ]; then \
	  sed "s/^  version: .*/  version: $$VERSION/" .speakeasy/gen.yaml > .speakeasy/gen.yaml.tmp && mv .speakeasy/gen.yaml.tmp .speakeasy/gen.yaml; \
	  echo "  Set terraform.version to $$VERSION"; \
	fi
```

Note the prerequisite on line 1 also changes from `overlay.yaml` to `overlay-v1.yaml`, and `--match 'v1.*'` is new.

The filter also fails safe. Verified against the current repo: `--match 'v0.*'` returns `v0.10.4` (57 matching tags), while `--match 'v1.*'` exits non-zero with "No names found" because no `v1.*` tag exists yet. The existing `2>/dev/null` plus the `if [ -n "$$VERSION" ]` guard turn that into "leave `gen.yaml` alone" rather than "silently fall back to 0.10.4" — which is exactly why Task 3 Step 2 must pass `VERSION=1.0.0-alpha.1` explicitly.

- [ ] **Step 6: Apply the mirror-image guard on `master`**

The only step in this plan that touches `master`. Same one-line idea, opposite tag series — done here rather than in a separate task because it is the same fix and shares the rationale above.

```bash
git checkout master
```

In `Makefile:6`, add the match filter so master only ever considers `v0.*` tags:

```makefile
	@VERSION=$$(git describe --tags --abbrev=0 --match 'v0.*' 2>/dev/null | sed 's/^v//'); \
```

Without this, master's version is correct only because `v1.*` tags are unreachable from master. If the review PR from Task 5 Step 3 is ever merged — or anyone cherry-picks from `v1` — `git describe` on master starts returning `1.0.0-alpha.1` and `make provider` silently stamps V0 code with a V1 version. You would not find out until release.

```bash
git add Makefile
git commit -m "Pin master version derivation to the v0 tag series

Guards against a v1.* tag becoming reachable from master and silently
stamping V0 code with a V1 version."
git checkout v1
```

- [ ] **Step 7: Create the tfplugindocs index template**

This must exist before Task 3 generates. `make provider` ends with `go generate ./...` (`Makefile:19`), which renders `docs/` — so if the template is missing at that point, you get a bannerless index and have to generate twice.

There is no `templates/` directory today, so the index page currently renders from tfplugindocs' default template plus the provider description. A one-line description is not a loud enough warning for the page the registry will show as "latest". Create `templates/index.md.tmpl`:

```markdown
---
page_title: "lumos Provider"
subcategory: ""
description: |-
  Lumos V1 provider — ALPHA. Breaking changes may land without notice.
---

# lumos Provider (1.x — ALPHA)

> ## ⚠️ This is an alpha release
>
> The **1.x** line of this provider targets the **Lumos V1 API** and may
> introduce breaking changes without prior notice. Pin an exact version.
>
> **For production use, stay on the 0.10.x line:** `version = "~> 0.10"`
>
> The 1.x line intentionally exposes **only** V1 API resources. It does not
> include the app store, access policy, pre-approval rule, requestable
> permission, user, or group resources available in 0.10.x. Those arrive in
> later 1.x alphas as they land in the V1 API. A configuration cannot use
> 0.10.x and 1.x at the same time — Terraform loads one version per provider
> address. Use separate root modules if you need both.

## Example Usage

{{ tffile "examples/provider/provider.tf" }}

<!-- schema generated by tfplugindocs -->
## Schema

{{ .SchemaMarkdown | trimspace }}
```

- [ ] **Step 8: Verify the build config without generating**

```bash
make -n provider VERSION=1.0.0-alpha.1 | head -5
test -f templates/index.md.tmpl && echo "template present"
grep -c "fast-api-v1" .speakeasy/workflow.yaml
```

Expected: the recipe echoes with `VERSION="1.0.0-alpha.1"` present, `template present`, and `1`. `-n` is dry-run, so nothing is written.

- [ ] **Step 9: Commit**

```bash
git add overlay-v1.yaml openapi-v1.json .speakeasy/workflow.yaml Makefile templates/
git commit -m "v1: add V1-only overlay, spec, and generation inputs

Keeps the four in-scope V1 paths, annotates App CRUD and Integration reads
for Speakeasy, sets the ALPHA provider description, and adds the docs banner
template. Knowledge Context is deferred: its multipart file field cannot
round-trip through Read, and its DELETE declares no success response."
```

---

### Task 3: Generate and verify

One `make provider` invocation produces the SDK, the provider, the examples, and the docs. Everything after it in this task is verification of that single output. Deliverable: a compiling, V1-only provider with correctly rendered docs.

**Files:**
- Modify: `.speakeasy/gen.yaml` (version, then re-enable persistent edits), `.speakeasy/gen.lock`, `internal/provider/*`, `internal/sdk/*`, `examples/*`, `tests/client_test.go`, `README.md`
- Generated: `docs/index.md`, `docs/resources/app.md`, `docs/data-sources/*.md`

**Interfaces:**
- Consumes: `overlay-v1.yaml`, the wired `.speakeasy/workflow.yaml`, and `templates/index.md.tmpl` from Task 2.
- Produces: a compiling provider exposing resource `lumos_app` and data sources `lumos_apps`, `lumos_integration`, `lumos_integrations`, with a `docs/` tree carrying the alpha banner. Task 4 exercises it against the live API.

- [ ] **Step 1: Record the pre-generation V0 file inventory**

You need this to prove generation actually removed the V0 files rather than leaving them behind.

```bash
ls internal/provider/*.go | wc -l > /tmp/v0_provider_count.txt
find examples -type f | wc -l > /tmp/v0_examples_count.txt
cat /tmp/v0_provider_count.txt /tmp/v0_examples_count.txt
```

Expected: `51` and `27` (27 examples, since Task 1 deleted 2 of the original 29).

- [ ] **Step 2: Generate**

```bash
make provider VERSION=1.0.0-alpha.1
```

This runs `speakeasy run --skip-compile --skip-versioning`, `go mod tidy`, `go build ./internal/provider/...`, `go test ./internal/provider/...`, then `go generate ./...` for the docs.

Expected: ends with `✓ Successfully generated and compiled the provider and ran tests`. If it fails inside `speakeasy run`, read the error before retrying — a failure here most likely means an annotation in Task 2 named a schema Speakeasy could not resolve.

- [ ] **Step 3: Assert no V0 resources survived — the critical check**

This is the test for the persistent-edits risk flagged in Task 1.

```bash
python3 - <<'PY'
import pathlib, re, sys

V0 = ['accesspolicy', 'accesspolicies', 'appstoreapp', 'appstoreappsettings',
      'preapprovalrule', 'requestablepermission', 'requestablepermissions',
      'user', 'users', 'group', 'groups']

files = sorted(p.name for p in pathlib.Path('internal/provider').glob('*.go'))
resurrected = [f for f in files
               if any(f.startswith(v + '_') for v in V0)]

prov = pathlib.Path('internal/provider/provider.go').read_text()
block = prov[prov.index('func (p *LumosProvider) Resources'):]
registered = sorted(set(re.findall(r'New(\w+?)(?:Resource|DataSource)\b', block)))

print('provider .go files:', len(files))
print('registered types:', registered)
print()
if resurrected:
    print('FAIL: V0 files present:', resurrected); sys.exit(1)
lower = [r.lower() for r in registered]
leaked = [r for r in lower if any(v in r for v in V0)]
if leaked:
    print('FAIL: V0 types registered:', leaked); sys.exit(1)
print('PASS: no V0 resources or data sources present')
PY
```

Expected: `PASS`, with `registered types` containing only App/Apps/Integration/Integrations variants.

**If this FAILs**, the persistent-edits reset in Task 1 Step 5 did not fully take. Recovery: `git checkout -- .`, re-verify `grep -c "persistentEdits\|trackedFiles" .speakeasy/gen.lock` returns 0 and `gen.yaml` has `enabled: "false"`, then re-run Step 2.

- [ ] **Step 4: Assert no V0 endpoint strings remain in the SDK**

```bash
grep -rn "appstore\|requestable_permission\|access_polic\|pre_approval\|/users\|/groups" internal/sdk/ --include=*.go | grep -v "_test.go" | head
```

Expected: no output. Any hit means an annotated V0 operation slipped through the overlay.

- [ ] **Step 5: Assert version and provider address**

```bash
grep "^  version:" .speakeasy/gen.yaml
grep "releaseVersion" .speakeasy/gen.lock
grep "Address:" main.go
```

Expected: `  version: 1.0.0-alpha.1`, `releaseVersion: 1.0.0-alpha.1`, and `Address: "registry.terraform.io/teamlumos/lumos"` unchanged per Global Constraints. If `gen.lock` still says `0.10.4`, the Makefile override did not apply.

- [ ] **Step 6: Assert the docs banner rendered and no V0 docs exist**

`make provider` already ran `go generate ./...`, so `docs/` is populated. Do not regenerate.

```bash
python3 - <<'PY'
import pathlib, sys
fail = []
idx = pathlib.Path('docs/index.md')
if not idx.exists():
    fail.append('docs/index.md was not generated')
else:
    t = idx.read_text()
    for needle in ['ALPHA', 'without prior notice', '~> 0.10']:
        if needle not in t:
            fail.append(f'docs/index.md missing {needle!r}')

V0 = ['access_polic', 'app_store', 'pre_approval', 'requestable', 'user', 'group']
stale = [str(p) for p in pathlib.Path('docs').rglob('*.md')
         if any(v in p.name for v in V0)]
if stale:
    fail.append(f'V0 docs present: {stale}')

for p in sorted(pathlib.Path('docs').rglob('*.md')):
    print(' ', p)
print()
if fail:
    print('FAIL'); [print(' -', f) for f in fail]; sys.exit(1)
print('PASS: alpha banner present, no V0 docs')
PY
```

Expected: `PASS`, listing roughly `docs/index.md`, `docs/resources/app.md`, `docs/data-sources/apps.md`, `docs/data-sources/integration.md`, `docs/data-sources/integrations.md`.

If the banner is missing, `templates/index.md.tmpl` was not created in Task 2 Step 7 — create it and re-run `go generate ./...` alone.

- [ ] **Step 7: Rewrite `tests/client_test.go` against V1 — required, this is a compile break**

`make provider` only builds `./internal/provider/...`, so Step 2 will not surface this. The existing file calls four symbols that no longer exist (`client.AppStore.GetAppStoreApp`, `operations.GetAppStoreAppRequest`, `client.Core.CreateApp`, `shared.AppInputCreate`) on two sub-clients that no longer exist (`Core`, `AppStore`).

The V1 replacements are already generated in the current tree, so these signatures are exact, not guesses:

| Call | Request type |
|---|---|
| `client.Apps.ListAppsV1` | `operations.ListAppsV1Request` |
| `client.Apps.GetAppV1` | `operations.GetAppV1Request` |
| `client.Apps.ConnectApp` | `shared.ConnectIntegrationInput` (bare struct, not a wrapper) |
| `client.Apps.ReconnectApp` | `operations.ReconnectAppRequest` |
| `client.Apps.DisconnectApp` | `operations.DisconnectAppRequest` |
| `client.Integrations.ListIntegrations` | `operations.ListIntegrationsRequest` |
| `client.Integrations.GetIntegration` | `operations.GetIntegrationRequest` |

One subtlety: `GetAppV1Request` currently has field `AppUUID`. The `x-speakeasy-name-override: "id"` added in Task 2 Step 1 renames it to **`ID`** — confirmed by the V0 precedent, where `GetAppRequest.ID` carries the tag `pathParam:"...,name=app_id"` while un-overridden `GetIntegrationRequest.AppID` keeps the raw name. Write `ID`, and if the compiler says `unknown field ID`, the override did not apply and Task 2 needs revisiting.

Replace `tests/client_test.go` lines 30-55 (both `TestGetApp` and `TestCreateApp`) with:

```go
func TestListIntegrations(t *testing.T) {
	client, ctx := GetNewClient(t)
	response, err := client.Integrations.ListIntegrations(ctx, operations.ListIntegrationsRequest{})
	if err != nil || response == nil {
		t.Fatalf(`client.Integrations.ListIntegrations() was unsuccessful, %v, %v`, response, err)
	}
	if response.PageIntegrationCatalogOutput == nil {
		t.Fatalf(`client.Integrations.ListIntegrations() returned no page payload, %v`, response)
	}
	if len(response.PageIntegrationCatalogOutput.Items) == 0 {
		t.Fatalf(`client.Integrations.ListIntegrations() returned zero integrations`)
	}
}

func TestGetAppV1(t *testing.T) {
	client, ctx := GetNewClient(t)
	appID := os.Getenv("APP_ID")
	if appID == "" {
		t.Skip(`APP_ID is not set`)
	}
	response, err := client.Apps.GetAppV1(ctx, operations.GetAppV1Request{
		ID: appID,
	})
	if err != nil || response == nil {
		t.Fatalf(`client.Apps.GetAppV1(%s) was unsuccessful, %v, %v`, appID, response, err)
	}
	if response.AppWithCustomAttributes == nil {
		t.Fatalf(`client.Apps.GetAppV1(%s) returned no app, %v`, appID, response)
	}
	if response.AppWithCustomAttributes.ID != appID {
		t.Fatalf(`client.Apps.GetAppV1(%s) returned id %q`, appID, response.AppWithCustomAttributes.ID)
	}
}
```

The `TestCreateApp` equivalent is deliberately dropped rather than ported: `ConnectApp` provisions a real integration against a live tenant, which is not something an unattended `go test` should do. Connect/reconnect/disconnect are exercised manually in Task 4 instead.

The `shared` import is still needed by `GetNewClient` (for `shared.Security`), so leave the import block alone; only the two test functions change.

- [ ] **Step 8: Full build and vet**

```bash
go build ./... && go vet ./... && go test ./internal/provider/...
```

Expected: all pass with no output from `build` and `vet`.

Then confirm the `tests/` package compiles. It hits the live API, so compile it without running it:

```bash
go vet ./tests/
```

Expected: no output. `unknown field ID in struct literal` means the Task 2 name override did not apply; `undefined: client.AppStore` means Step 7 was not completed.

- [ ] **Step 9: Update the README**

`README.md` currently documents V0 usage — `source = "teamlumos/lumos"` with V0 examples at lines 88, 227, and 322. Put the alpha block at the very top, immediately after the title, and replace the V0 resource examples with the `lumos_app` example from `examples/resources/lumos_app/resource.tf`. Use the exact Global Constraints wording:

```markdown
> ## ⚠️ ALPHA — 1.x targets the Lumos V1 API
>
> The 1.x line may introduce breaking changes without prior notice. Pin an
> exact version: `version = "1.0.0-alpha.1"`.
>
> **For production, stay on 0.10.x:** `version = "~> 0.10"`
```

- [ ] **Step 10: Re-enable persistent edits and re-run to establish the new baseline**

Now that a clean V1 tree exists, turn the feature back on so future regenerations preserve manual edits.

```bash
python3 - <<'EOF'
p = '.speakeasy/gen.yaml'
s = open(p).read()
s = s.replace('  persistentEdits:\n    enabled: "false"', '  persistentEdits:\n    enabled: "true"')
open(p, 'w').write(s)
EOF
git add -A && git commit -m "v1: first V1-only generation at 1.0.0-alpha.1

Generates the V1 provider (lumos_app resource, apps/integration/integrations
data sources), regenerates docs with the alpha banner, and ports the SDK
smoke tests from V0 AppStore/Core calls to V1 Apps/Integrations."
make provider VERSION=1.0.0-alpha.1
grep -A3 "persistentEdits:" .speakeasy/gen.lock
```

Expected: a fresh `pristine_commit_hash` / `pristine_tree_hash` pointing at the commit you just made — **not** `23d14bec32f4c6f6341df608eda8e48526d0c4f1`.

- [ ] **Step 11: Commit**

```bash
git add -A
git commit -m "v1: re-enable persistent edits with V1 baseline"
```

---

### Task 4: Verify against the live V1 API

Generation succeeding proves nothing about whether the resource actually functions. This is the functional gate, and it is deliberately separate from Task 5 because publishing a tag is not cleanly reversible.

**Files:**
- Create: `/tmp/lumos-v1-test/main.tf`, `/tmp/lumos-v0-check/main.tf` (throwaway, not committed)

**Interfaces:**
- Consumes: the built provider from Task 3.
- Produces: confirmed evidence that `lumos_app` and the integration data sources work end-to-end. Task 5 releases only after this passes.

- [ ] **Step 1: Build and install the provider into a local filesystem mirror**

```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/teamlumos/lumos/1.0.0-alpha.1/darwin_arm64
go build -o ~/.terraform.d/plugins/registry.terraform.io/teamlumos/lumos/1.0.0-alpha.1/darwin_arm64/terraform-provider-lumos_v1.0.0-alpha.1 .
ls -la ~/.terraform.d/plugins/registry.terraform.io/teamlumos/lumos/1.0.0-alpha.1/darwin_arm64/
```

This is the *unpacked* filesystem-mirror layout (`HOSTNAME/NAMESPACE/TYPE/VERSION/TARGET/`), which is the right shape for a locally built binary. Adjust `darwin_arm64` if you are not on Apple silicon.

- [ ] **Step 2: Write a test configuration**

Create `/tmp/lumos-v1-test/main.tf`:

```hcl
terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "1.0.0-alpha.1"
    }
  }
}

provider "lumos" {}

data "lumos_integrations" "all" {}

output "integration_count" {
  value = length(data.lumos_integrations.all.items)
}

output "first_integration" {
  value = try(data.lumos_integrations.all.items[0], null)
}
```

- [ ] **Step 3: Confirm the version pin resolves to the local build**

```bash
cd /tmp/lumos-v1-test
export LUMOS_ACCESS_TOKEN="<a real token for a test tenant>"
terraform init
```

Expected: `- Installing teamlumos/lumos v1.0.0-alpha.1...` and `- Installed teamlumos/lumos v1.0.0-alpha.1 (unauthenticated)`. The `(unauthenticated)` note is expected for a local mirror build — it means no signature check, not a failure.

If it instead tries to reach the registry and fails, the mirror path or platform directory is wrong.

- [ ] **Step 4: Exercise the read path against the live API**

```bash
terraform plan
terraform apply -auto-approve
terraform output integration_count
```

Expected: a non-zero `integration_count` and a populated `first_integration`. This proves auth, the `x-speakeasy-retries` config, pagination mapping on `Page_IntegrationCatalogOutput_`, and the data source schema all work.

- [ ] **Step 5: Verify the V0 provider still resolves independently**

This is the regression test for the whole "reuse the address" decision.

```bash
mkdir -p /tmp/lumos-v0-check && cd /tmp/lumos-v0-check
cat > main.tf <<'EOF'
terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "~> 0.10"
    }
  }
}
EOF
terraform init
```

Expected: `Installing teamlumos/lumos v0.10.4` from the public registry. This confirms `~> 0.10` does **not** pick up the 1.x alpha, which is the property protecting existing customers.

- [ ] **Step 6: Record the results**

Note for the PR description: integration count returned, whether `terraform apply` was clean, and the resolved version in each of the two directories. If `lumos_app` create/update/delete was also exercised against a test tenant, say which integration was used.

Nothing to commit — this is a verification gate.

---

### Task 5: Release `1.0.0-alpha.1` and verify ingestion

**Files:**
- Modify: `.goreleaser.yml` (one line, `v1` branch only)

**Interfaces:**
- Consumes: the verified `v1` branch from Task 4.
- Produces: `teamlumos/lumos` `1.0.0-alpha.1` installable from the public registry by exact pin, with the V0 line proven still releasable.

- [ ] **Step 1: Set the GitHub release to non-prerelease**

GoReleaser defaults `release.prerelease` to `auto`, which detects the `-alpha` suffix and marks the **GitHub release** as a prerelease. HashiCorp's docs state the registry creates a webhook "subscribed to release events" but do not say whether GitHub prereleases are ingested. Rather than gamble on it, make it an ordinary GitHub release carrying a prerelease *semver tag* — the tag is what determines Terraform's prerelease behavior, not the GitHub flag.

Add to the `release:` block in `.goreleaser.yml` (currently at lines 61-65):

```yaml
release:
  prerelease: false
  extra_files:
    - glob: "terraform-registry-manifest.json"
      name_template: "{{ .ProjectName }}_{{ .Version }}_manifest.json"
```

- [ ] **Step 2: Verify goreleaser produces correct artifacts without publishing**

```bash
goreleaser release --snapshot --clean --skip=sign
ls dist/*.zip | head
```

Expected: zips named `terraform-provider-lumos_1.0.0-alpha.1-SNAPSHOT-<sha>_<os>_<arch>.zip`. The `-SNAPSHOT` suffix is a snapshot-mode artifact; the real tagged build drops it. Confirm the version segment reads `1.0.0-alpha.1`, not `0.10.4`.

- [ ] **Step 3: Push the branch and open a review PR**

```bash
git add .goreleaser.yml
git commit -m "v1: publish tagged alphas as normal GitHub releases"
git push -u origin v1
gh pr create --base master --head v1 --draft \
  --title "V1-only provider (1.0.0-alpha.1)" \
  --body "Long-lived v1 branch. NOT for merge to master — this PR is for review only.

Generates a V1-API-only provider published as teamlumos/lumos 1.0.0-alpha.1.
See plans/2026-08-20-v1-alpha-provider.md for the full design.

Knowledge Context is deferred: multipart file field cannot round-trip through
Read, and DELETE /v1/knowledge/context/{id} declares no success response."
```

**Do not merge this PR.** The `v1` branch is long-lived; the PR exists for review, and is opened as a draft to reduce the risk of an accidental merge. Merging it would make `v1.0.0-alpha.1` reachable from master — which the Task 2 Step 6 guard now defends against, but the branch content would still be wrong.

- [ ] **Step 4: Tag and release**

```bash
git tag v1.0.0-alpha.1
git push origin v1.0.0-alpha.1
gh run watch
```

Expected: `terraform_publish.yaml` fires (it triggers on `v*` tags from any branch, so it needs no changes) and GoReleaser publishes a signed GitHub release.

- [ ] **Step 5: Verify the registry ingested the prerelease — the open question**

Wait a few minutes, then:

```bash
curl -s https://registry.terraform.io/v1/providers/teamlumos/lumos/versions \
  | python3 -c "import json,sys; d=json.load(sys.stdin); vs=[v['version'] for v in d['versions']]; print('1.0.0-alpha.1 present:', '1.0.0-alpha.1' in vs)"
```

Expected: `1.0.0-alpha.1 present: True`.

**If False**, the webhook did not ingest it. Recovery, in order: (1) go to the provider settings page on registry.terraform.io and click **Resync**; (2) confirm the GitHub release is not marked prerelease or draft; (3) confirm the release contains the `_SHA256SUMS`, `_SHA256SUMS.sig`, and `_manifest.json` assets.

- [ ] **Step 6: Verify a clean install from the public registry**

```bash
rm -rf /tmp/lumos-v1-registry && mkdir -p /tmp/lumos-v1-registry && cd /tmp/lumos-v1-registry
cat > main.tf <<'EOF'
terraform {
  required_providers {
    lumos = {
      source  = "teamlumos/lumos"
      version = "1.0.0-alpha.1"
    }
  }
}
EOF
TF_CLI_CONFIG_FILE=/dev/null terraform init
```

Expected: `Installing teamlumos/lumos v1.0.0-alpha.1` with a signature check that passes (no `(unauthenticated)` note this time — it is GPG-signed and served by the registry). `TF_CLI_CONFIG_FILE=/dev/null` bypasses your local mirror so this genuinely tests the registry.

- [ ] **Step 7: Confirm the accepted cost, so it is not a surprise**

```bash
open "https://registry.terraform.io/providers/teamlumos/lumos/latest"
```

Expected, and **accepted per Design Decision 4**: the page shows `1.0.0-alpha.1` as latest, with the alpha banner from `templates/index.md.tmpl` rendered at the top. Confirm the banner is actually visible above the "USE PROVIDER" snippet — that is the whole mitigation. If it did not render, fix the template and ship `1.0.0-alpha.2`.

- [ ] **Step 8: Prove the V0 line is still releasable**

Now that a `v1.*` tag exists, verify the Task 2 Step 6 guard actually holds. This is the test for the whole "does this break V0" question.

```bash
cd /Users/grahamseyffert/develop/terraform-provider-lumos
git checkout master
git describe --tags --abbrev=0 --match 'v0.*' | sed 's/^v//'
```

Expected: `0.10.4` — **not** `1.0.0-alpha.1`.

Then prove a V0 release still builds end to end. Do not tag a throwaway version; snapshot mode exercises the whole build and archive path without publishing:

```bash
goreleaser release --snapshot --clean --skip=sign
ls dist/*.zip | head -3
git checkout v1
```

Expected: zips whose version segment reads `0.10.4-SNAPSHOT-<sha>`, proving master builds a V0 artifact at a V0 version with `.goreleaser.yml` untouched.

**What extending V0 looks like afterwards.** Unchanged from today: edit `overlay.yaml` on master, `make provider`, verify, `git tag v0.11.0`, push. `terraform_publish.yaml` fires on any `v*` tag from any branch, so it needs no changes and never has.

**One consequence worth telling the account team.** A new V0 release will not appear as "latest" on the registry page while `1.0.0-alpha.1` exists, because `1.0.0-alpha.1` outranks every `0.x` under SemVer precedence and the registry's sort is [strictly SemVer compliant](https://developer.hashicorp.com/terraform/registry/faq). A legacy customer on `~> 0.10` still *resolves* 0.11.0 correctly — `~>` ignores prereleases, so nothing breaks — but they will not discover it from the page's default view. If that matters for a specific legacy rollout, announce the release directly rather than relying on the registry page.

---

## Self-Review

**Spec coverage:** Every design decision maps to a task. Same-address/1.x → Task 3 Step 5 and Task 4 Step 5. V0/V1 non-coexistence → documented in the Task 2 Step 7 banner, verified in Task 4 Step 5. Public prerelease → Task 5. Accepted "latest" cost → mitigated in Task 2 Step 7, confirmed in Task 5 Step 7. The `.genignore` assessment that was explicitly requested → Task 1 Steps 2-4, with the finding that all three entries are obsolete and why `requestable_permissions_hook.go` compiles but is dead. The persistent-edits risk → Task 1 Step 5, verified in Task 3 Step 3 with a stated recovery path. V0 remaining releasable → Task 2 Step 6, proven in Task 5 Step 8.

**Placeholder scan:** No TBDs. Every code step carries runnable content. The one genuine unknown — whether the registry ingests a GitHub prerelease — is isolated to Task 5 Step 5 with a concrete pass/fail check and a three-step recovery, rather than left as an assumption. The `<a real token for a test tenant>` in Task 4 Step 3 is a required human input, not a placeholder for missing plan content.

**Type consistency:** Entity names are consistent across Task 2 (`App`, `Apps`, `Integration`, `Integrations`) and the Task 3 Step 3 assertion. Schema names (`AppWithCustomAttributes`, `ConnectIntegrationInput`, `Page_AppWithCustomAttributes_`, `IntegrationCatalogOutput`, `Page_IntegrationCatalogOutput_`) were read from `openapi.json` and match. The `app_uuid` → `id` override is applied to all three methods on `/v1/apps/{app_uuid}`, and the plan notes explicitly that this is distinct from `ConnectIntegrationInput.app_id`.

Every Go identifier in the Task 3 Step 7 test code was verified against the currently generated SDK, not inferred: `client.Apps.GetAppV1` / `client.Integrations.ListIntegrations` (`internal/sdk/apps.go:596`, `internal/sdk/integrations.go:37`), `ListIntegrationsResponse.PageIntegrationCatalogOutput` (`internal/sdk/models/operations/listintegrations.go:52`), `PageIntegrationCatalogOutput.Items` (`internal/sdk/models/shared/pageintegrationcatalogoutput.go:7`), and `AppWithCustomAttributes.ID` (`internal/sdk/models/shared/appwithcustomattributes.go:8`). The one field the plan expects to *change* — `GetAppV1Request.AppUUID` → `ID` — is called out with the V0 precedent that establishes the behavior and a specific compiler error to watch for.

**Consolidation from 8 tasks to 5:**
- **Old Task 4 (docs) folded into Tasks 2 and 3.** This fixed a real ordering defect, not just task count: `make provider` already ends with `go generate ./...` (`Makefile:19`), so old Task 3 generated docs and old Task 4 immediately regenerated them — and because the banner template did not exist until old Task 4, the first generation would have produced a bannerless `docs/index.md`. Template creation now precedes generation (Task 2 Step 7) and docs verification sits with the generation it verifies (Task 3 Step 6).
- **Old Task 8 (V0 guard) split by timing rather than kept whole.** The `--match` edit is one idea applied to two branches, so both edits now live together in Task 2 Steps 5-6 instead of restating the rationale twice. Its verification genuinely could not move there — it requires a `v1.*` tag to exist — so it lives in Task 5 Step 8 alongside the other post-release checks.
- **Old Task 7 (Knowledge Context) demoted out of the task list.** It had no steps and no deliverable; it was deferred scope wearing a task number. It is now a "Deferred Scope" section with the blockers and a concrete pick-up procedure.
- **Deliberately NOT merged:** old Tasks 5 and 6 (now 4 and 5) stay separate. Task 4 is reversible local verification; Task 5 publishes a tag that cannot be cleanly unpublished. That boundary is the most valuable gate in the plan. Tasks 1 and 2 also stay separate — combined they would run ~16 steps across two unrelated failure modes (baseline reset vs. annotation correctness), and keeping them apart means a Task 3 generation failure is bisectable.

**Corrections applied during earlier review:**
- Task 3 Step 1 expected provider-file count corrected from 54 to 51 (measured).
- `tests/client_test.go` identified as a hard compile break, not optional cleanup. It calls `client.AppStore.GetAppStoreApp` and `client.Core.CreateApp` on sub-clients that cease to exist under V1-only generation, and `make provider` cannot catch it because that target only builds `./internal/provider/...`.

**Verified rather than assumed:** `git describe --tags --abbrev=0 --match 'v0.*'` returns `v0.10.4` in the current repo (57 matching tags), and `--match 'v1.*'` exits non-zero rather than falling back. The six V1 paths and all five annotated schemas are present in the tracked `original-openapi.json`, so Task 2's spec derivation has a real foundation. Neither `openapi.json` nor `openapi-v1.json` is covered by `.gitignore`, so Task 2 Step 9's `git add` will not silently no-op.
