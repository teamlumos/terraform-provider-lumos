package listplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.List = ListEmptyWhenNullPlanModifier{}

// ListEmptyWhenNullPlanModifier plans an empty list for a list the config
// leaves out, so omitting it clears it and the plan shows that. It stands in
// for a spec default of [], which the generator doesn't turn into a
// Terraform default for a list of objects or a list inside an object.
type ListEmptyWhenNullPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v ListEmptyWhenNullPlanModifier) Description(_ context.Context) string {
	return "Plans an empty list when the config leaves the list out."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v ListEmptyWhenNullPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyList performs the plan modification.
func (v ListEmptyWhenNullPlanModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	resp.PlanValue = types.ListValueMust(req.ConfigValue.ElementType(ctx), []attr.Value{})
}

func EmptyWhenNull() planmodifier.List {
	return ListEmptyWhenNullPlanModifier{}
}
