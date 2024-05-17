import json

def add_plan_modifiers_recursively(attribute, key):
    key = "list" if "list" in key else key
    plan_modifiers_block = {
        "plan_modifiers": [
            {
                "custom": {
                    "imports": [
                        {
                            "path": f"github.com/hashicorp/terraform-plugin-framework/resource/schema/{key}planmodifier"
                        }
                    ],
                    "schema_definition": f"{key}planmodifier.UseStateForUnknown()"
                }
            }
        ]
    }
    
    if isinstance(attribute, dict):
        if "computed_optional_required" in attribute and "computed" in attribute["computed_optional_required"]:
            attribute.update(plan_modifiers_block)
        
        for key, value in attribute.items():
            if isinstance(value, (dict, list)):
                add_plan_modifiers_recursively(value, key)

            if key == "nested_object":
                for attr in value["attributes"]:
                    add_plan_modifiers_recursively(attr, key)
    elif isinstance(attribute, list):
        for item in attribute:
            add_plan_modifiers_recursively(item, "list")



with open("tmp-provider-code-spec.json", "r") as f:
    data = json.load(f)
    for resource in data["resources"]:
        for attribute in resource["schema"]["attributes"]:
            key = "string" if "string" in attribute.keys() else "list" 
            add_plan_modifiers_recursively(attribute, key)

print(json.dumps(data, indent=2))
