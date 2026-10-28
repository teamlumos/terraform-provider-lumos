package boolplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

var _ planmodifier.Bool = BoolUnknownOnReconnectPlanModifier{}

// BoolUnknownOnReconnectPlanModifier plans a lumos_app attribute that only
// a reconnect changes. An update that changes the credentials reconnects,
// so the value stays unknown until the API answers; any other update keeps
// the current value. On disconnected, it also plans false for a
// disconnected integration, so the plan reconnects it.
type BoolUnknownOnReconnectPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v BoolUnknownOnReconnectPlanModifier) Description(_ context.Context) string {
	return "Keeps the current value unless the update reconnects the integration."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v BoolUnknownOnReconnectPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyBool performs the plan modification.
func (v BoolUnknownOnReconnectPlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	if req.Path.Equal(path.Root("disconnected")) {
		plansReconnect, diags := utils.AppPlansReconnect(ctx, req.Plan, req.State)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		if plansReconnect {
			resp.PlanValue = types.BoolValue(false)
			return
		}
	}
	if !req.PlanValue.IsUnknown() {
		return
	}
	reconnects, diags := utils.AppReconnectsOnUpdate(ctx, req.Plan, req.State)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() || reconnects {
		return
	}
	resp.PlanValue = req.StateValue
}

func UnknownOnReconnect() planmodifier.Bool {
	return BoolUnknownOnReconnectPlanModifier{}
}
