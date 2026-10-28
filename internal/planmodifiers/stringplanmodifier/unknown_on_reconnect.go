package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

var _ planmodifier.String = StringUnknownOnReconnectPlanModifier{}

// StringUnknownOnReconnectPlanModifier plans a lumos_app attribute that only
// a reconnect changes. An update that changes the credentials reconnects,
// so the value stays unknown until the API answers; any other update keeps
// the current value. A plan that reconnects a disconnected integration
// leaves it unknown too.
type StringUnknownOnReconnectPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringUnknownOnReconnectPlanModifier) Description(_ context.Context) string {
	return "Keeps the current value unless the update reconnects the integration."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringUnknownOnReconnectPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyString performs the plan modification.
func (v StringUnknownOnReconnectPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		plansReconnect, diags := utils.AppPlansReconnect(ctx, req.Plan, req.State)
		resp.Diagnostics.Append(diags...)
		if !diags.HasError() && plansReconnect {
			resp.PlanValue = types.StringUnknown()
		}
		return
	}
	reconnects, diags := utils.AppReconnectsOnUpdate(ctx, req.Plan, req.State)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() || reconnects {
		return
	}
	resp.PlanValue = req.StateValue
}

func UnknownOnReconnect() planmodifier.String {
	return StringUnknownOnReconnectPlanModifier{}
}
