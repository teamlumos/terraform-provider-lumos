import argparse
from dataclasses import dataclass
import requests
import sys
import json
import urllib

REQUESTABLE_PERMISSION_RESOURCE = "lumos-appstore_requestable_permission"
USER_RESOURCE = "lumos-appstore_user"


class LumosAPI:
    LUMOS_API_ENDPOINT = "https://api.lumos.com"
    LUMOS_API_TOKEN = "LUMOS_API_KEY"

    _headers = None

    def __init__(self) -> None:
        self._headers = {
            "Authorization": f"Bearer {self.LUMOS_API_TOKEN}",
            "Content-Type": "application/json",
        }

    def search_user(self, email: str) -> dict:
        response = requests.get(
            f"{self.LUMOS_API_ENDPOINT}/users?search_term={urllib.parse.quote(email)}",
            headers=self._headers,
        )
        if response.status_code != 200:
            print(
                f"Failed to get user with email {email}. "
                f"Got {response.status_code} response: {response.text}"
            )

        response_json = response.json()
        if response_json["total"] == 0:
            return None
        if response_json["total"] > 1:
            print(
                f"Found more than one user with email {email}. "
                f"Returning the first one."
            )
        return response_json["items"][0]


    def get_user(self, user_id: str) -> dict:
        response = requests.get(
            f"{self.LUMOS_API_ENDPOINT}/users/{user_id}",
            headers=self._headers,
        )
        if response.status_code != 200:
            print(
                f"Failed to get user with ID {user_id}. "
                f"Got {response.status_code} response: {response.text}"
            )

        return response.json()

    def get_requestable_permission(self, permission_id: str) -> dict:
        response = requests.get(
            f"{self.LUMOS_API_ENDPOINT}/appstore/requestable_permissions/{permission_id}",
            headers=self._headers,
        )
        if response.status_code != 200:
            print(
                f"Failed to get requestable permissions with ID {permission_id}. "
                f"Got {response.status_code} response: {response.text}"
            )

        return response.json()

@dataclass
class User:
    id: str
    given_name: str
    family_name: str
    email: str

    @property
    def resource_name(self) -> str:
        

@dataclass
class Permission:
    app_id: str
    label: str
    visible_in_appstore: str
    allowed_groups: list[str]
    manager_approval_required: bool
    require_additional_approval: bool
    approver_groups_stage_1: list[str]
    approver_users_stage_1: list[str]
    approver_groups_stage_2: list[str]
    approver_users_stage_2: list[str]
    access_removal_inline_webhook: str
    request_validation_inline_webhook: str
    provisioning_inline_webhook: str
    manual_steps_needed: bool
    manual_instructions: str
    time_based_access_options: list[str]
    provisioning_group: str

    @property
    def resource_name(self) -> str:
        snake_case_string = self.label.replace(' ', '_')
        snake_case_string = snake_case_string.lower()

        return ''.join(letter for letter in snake_case_string if letter.isalnum() or letter == "_")

def parse_user_response(user_json: dict) -> User:
    return User(
        id=user_json["id"],
        given_name=user_json["given_name"],
        family_name=user_json["family_name"],
        email=user_json["email"],
    )

def parse_permission_response(permission_json: dict) -> Permission:
    allowed_group_ids: list[str] = []
    for group in permission_json["request_config"].get("allowed_groups", {}).get("groups", []):
        allowed_group_ids.append(group["id"])

    approvers_stage_1: list[str] = []
    approvers_stage_2: list[str] = []

    approvers = permission_json["request_config"].get("request_approval_config", {}).get("request_approval_staged", [])

    if len(approvers) > 0:
        approvers_stage_1 = approvers[0].get("approvers")

    if len(approvers) > 1:
        approvers_stage_2 = approvers[1].get("approvers")

    user_ids_approvers_stage_1: list[str] = []
    group_ids_approvers_stage_1: list[str] = []
    if len(approvers_stage_1) > 0:
        for a in approvers_stage_1:
            if a["type"] == "USER":
                user_ids_approvers_stage_1.append(a["id"])
            else:
                group_ids_approvers_stage_1.append(a["id"])

    user_ids_approvers_stage_2: list[str] = []
    group_ids_approvers_stage_2: list[str] = []
    if len(approvers_stage_2) > 0:
        for a in approvers_stage_2:
            if a["type"] == "USER":
                user_ids_approvers_stage_2.append(a["id"])
            else:
                group_ids_approvers_stage_2.append(a["id"])

    try:
        access_removal_webhoook_id = permission_json["request_config"].get("access_removal_inline_webhook", {}).get("id")
    except Exception:
        access_removal_webhoook_id = None

    try:
        request_validation_webhook_id = permission_json["request_config"].get("request_validation_inline_webhook", {}).get("id")
    except Exception:
        request_validation_webhook_id = None

    try:
        provisioning_inline_webhook = permission_json["request_config"].get("request_fulfillment_config", {}).get("provisioning_inline_webhook", {}).get("id")
    except Exception:
        provisioning_inline_webhook = None

    try:
        provisioning_group = permission_json["request_config"].get("request_fulfillment_config", {}).get("provisioning_group", {}).get("id")
    except Exception:
        provisioning_group = None

    if permission_json["request_config"].get("request_fulfillment_config", {}).get("manual_instructions"):
        manual_instructions = permission_json["request_config"].get("request_fulfillment_config", {}).get("manual_instructions")
    else:
        manual_instructions = None

    return Permission(
        app_id=permission_json["app_id"],
        label=permission_json["label"],
        visible_in_appstore=True if permission_json["request_config"]["appstore_visibility"] == "VISIBLE" else False,
        allowed_groups=allowed_group_ids,
        manager_approval_required=True if permission_json["request_config"].get("request_approval_config", {}).get("manager_approval") == "INITIAL_APPROVAL" else False,
        require_additional_approval=str(permission_json["request_config"].get("request_approval_config", {}).get("require_additional_approval", False)).lower(),
        approver_users_stage_1=user_ids_approvers_stage_1,
        approver_groups_stage_1=group_ids_approvers_stage_1,
        approver_users_stage_2=user_ids_approvers_stage_2,
        approver_groups_stage_2=group_ids_approvers_stage_2,
        access_removal_inline_webhook=access_removal_webhoook_id,
        request_validation_inline_webhook=request_validation_webhook_id,
        provisioning_inline_webhook=provisioning_inline_webhook,
        provisioning_group=provisioning_group,
        manual_steps_needed=permission_json["request_config"].get("request_fulfillment_config", {}).get("manual_steps_needed"),
        manual_instructions=manual_instructions,
        time_based_access_options=permission_json["request_config"].get("request_fulfillment_config", {}).get("time_based_access")
    )

def search_user(email: str) -> User:
    user_json = LumosAPI().search_user(email)
    return parse_user_response(user_json)

def get_user(id: str) -> User:
    user_json = LumosAPI().get_user(id)
    return parse_user_response(user_json)

def get_user(id: str) -> User:
    user_json = LumosAPI().get_user(id)
    return parse_user_response(user_json)

def get_permission(id: str) -> Permission:
    permission_json = LumosAPI().get_requestable_permission(id)
    return parse_permission_response(permission_json)


def generate_resource_definition(permission_id: str):
    permission = get_permission(permission_id)

    access_removal_inline_webhook_value = f'"{permission.access_removal_inline_webhook}"' if permission.access_removal_inline_webhook is not None else "null"
    request_validation_inline_webhook_value = f'"{permission.request_validation_inline_webhook=}"' if permission.request_validation_inline_webhook is not None else "null"
    provisioning_inline_webhook_value = f'"{permission.provisioning_inline_webhook}"' if permission.provisioning_inline_webhook is not None else "null"
    provisioning_group_value = f'"{permission.provisioning_group}"' if permission.provisioning_group is not None else "null"
    manual_instructions_value = f'"{permission.manual_instructions}"' if permission.manual_instructions is not None else "null"

    print(f'resource "{REQUESTABLE_PERMISSION_RESOURCE}" "{permission.resource_name}"' + "{")
    print(f'\t \tapp_id="{permission.app_id}"')
    print(f'\t \tlabel="{permission.label}"')
    print(f'\t \tvisible_in_appstore={str(permission.visible_in_appstore).lower()}')
    print(f'\t \tallowed_groups={json.dumps(permission.allowed_groups)}')
    print(f'\t \tmanager_approval_required={str(permission.manager_approval_required).lower()}')
    print(f'\t \trequire_additional_approval={permission.require_additional_approval}')
    print(f'\t \tapprover_groups_stage_1={json.dumps(permission.approver_groups_stage_1)}')
    print(f'\t \tapprover_users_stage_1={json.dumps(permission.approver_users_stage_1)}')
    print(f'\t \tapprover_groups_stage_2={json.dumps(permission.approver_groups_stage_2)}')
    print(f'\t \tapprover_users_stage_2={json.dumps(permission.approver_users_stage_2)}')
    print(f'\t \taccess_removal_inline_webhook={access_removal_inline_webhook_value}')
    print(f'\t \trequest_validation_inline_webhook={request_validation_inline_webhook_value}')
    print(f'\t \tprovisioning_inline_webhook={provisioning_inline_webhook_value}')
    print(f'\t \tmanual_steps_needed={str(permission.manual_steps_needed).lower()}')
    print(f'\t \tmanual_instructions={manual_instructions_value}')
    print(f'\t \ttime_based_access_options={json.dumps(permission.time_based_access_options)}')
    print(f'\t \tprovisioning_group={provisioning_group_value}')
    print("}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Terraform Resource Generation for import from Lumos.")
    parser.add_argument("resource", help="The resource we want to import.")
    parser.add_argument("id", help="Stable ID from the resource we want to import.")

    args = parser.parse_args()

    if args.resource == REQUESTABLE_PERMISSION_RESOURCE:
        generate_resource_definition(args.id)
    else:
        sys.exit(f'Unknown resource {args.resource}')


"""
    Usage example: 
    python3 generate_lumos_terraform_resource_for_import.py lumos-appstore_requestable_permission 798a8e16-ad45-0f43-4fed-2c33175e8a77
"""
