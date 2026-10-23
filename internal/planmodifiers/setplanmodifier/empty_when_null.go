package setplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Set = SetEmptyWhenNullPlanModifier{}

// SetEmptyWhenNullPlanModifier plans an empty set for a set the config
// leaves out, so omitting it clears it and the plan shows that. It stands in
// for a spec default of [], which the generator only turns into a Terraform
// default when the response schema has it too, and lumos_app shows that
// schema read-only.
type SetEmptyWhenNullPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v SetEmptyWhenNullPlanModifier) Description(_ context.Context) string {
	return "Plans an empty set when the config leaves the set out."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v SetEmptyWhenNullPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifySet performs the plan modification.
func (v SetEmptyWhenNullPlanModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	resp.PlanValue = types.SetValueMust(req.ConfigValue.ElementType(ctx), []attr.Value{})
}

func EmptyWhenNull() planmodifier.Set {
	return SetEmptyWhenNullPlanModifier{}
}
