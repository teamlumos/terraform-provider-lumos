package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var _ planmodifier.Object = ObjectEmptyWhenNullPlanModifier{}

// ObjectEmptyWhenNullPlanModifier plans an empty object for an object the
// config leaves out: every list or set in it empty and everything else
// null. Omitting the object then clears each of its settings, and the plan
// shows that. It stands in for a spec default, which the generator doesn't
// turn into a Terraform default for an object.
type ObjectEmptyWhenNullPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v ObjectEmptyWhenNullPlanModifier) Description(_ context.Context) string {
	return "Plans empty lists and null values when the config leaves the object out."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v ObjectEmptyWhenNullPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyObject performs the plan modification.
func (v ObjectEmptyWhenNullPlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	attrTypes := req.ConfigValue.AttributeTypes(ctx)
	attrs := make(map[string]attr.Value, len(attrTypes))
	for name, attrType := range attrTypes {
		switch t := attrType.(type) {
		case types.ListType:
			attrs[name] = types.ListValueMust(t.ElemType, []attr.Value{})
		case types.SetType:
			attrs[name] = types.SetValueMust(t.ElemType, []attr.Value{})
		default:
			null, err := attrType.ValueFromTerraform(ctx, tftypes.NewValue(attrType.TerraformType(ctx), nil))
			if err != nil {
				resp.Diagnostics.AddAttributeError(req.Path, "Planning an omitted object", err.Error())
				return
			}
			attrs[name] = null
		}
	}
	empty, diags := types.ObjectValue(attrTypes, attrs)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	resp.PlanValue = empty
}

func EmptyWhenNull() planmodifier.Object {
	return ObjectEmptyWhenNullPlanModifier{}
}
