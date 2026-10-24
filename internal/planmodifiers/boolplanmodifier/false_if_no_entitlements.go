package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Bool = BoolFalseIfNoEntitlementsPlanModifier{}

// BoolFalseIfNoEntitlementsPlanModifier plans a pre-approval rule's
// request_has_entitlements when the config leaves it out. The API reads a
// rule with none of its three entitlements fields set back as
// request_has_entitlements = false, and leaves the field out when one of the
// two lists is set, so the plan is false or null to match.
type BoolFalseIfNoEntitlementsPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v BoolFalseIfNoEntitlementsPlanModifier) Description(_ context.Context) string {
	return "Plans false when the configuration sets no entitlements field, and null when it sets one of the entitlement lists."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v BoolFalseIfNoEntitlementsPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyBool performs the plan modification.
func (v BoolFalseIfNoEntitlementsPlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	for _, name := range []string{"request_includes_entitlements", "request_does_not_include_entitlements"} {
		var list types.Set
		diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName(name), &list)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		if !list.IsNull() {
			resp.PlanValue = types.BoolNull()
			return
		}
	}
	resp.PlanValue = types.BoolValue(false)
}

func FalseIfNoEntitlements() planmodifier.Bool {
	return BoolFalseIfNoEntitlementsPlanModifier{}
}
