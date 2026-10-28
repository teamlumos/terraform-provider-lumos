package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

var _ planmodifier.String = StringShowsOverridePlanModifier{}

// StringShowsOverridePlanModifier plans a lumos_app display field, such as
// label, that the API reads back as its <name>_override where one is set
// and as the app's default otherwise. A set override is what the field will
// show. Without one, the field keeps its current value, unless the update
// clears the override or reconnects the integration, by changing the
// credentials or because it's disconnected; either can change the default
// it shows.
type StringShowsOverridePlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringShowsOverridePlanModifier) Description(_ context.Context) string {
	return "Plans the override where one is set, and otherwise keeps the current value."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringShowsOverridePlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyString performs the plan modification.
func (v StringShowsOverridePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	overridePath := path.Root(req.Path.String() + "_override")

	var override types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, overridePath, &override)...)
	if resp.Diagnostics.HasError() || override.IsUnknown() {
		return
	}
	if !req.PlanValue.IsUnknown() {
		if !override.IsNull() || req.State.Raw.IsNull() {
			return
		}
		plansReconnect, diags := utils.AppPlansReconnect(ctx, req.Plan, req.State)
		resp.Diagnostics.Append(diags...)
		if !diags.HasError() && plansReconnect {
			resp.PlanValue = types.StringUnknown()
		}
		return
	}
	if !override.IsNull() {
		resp.PlanValue = override
		return
	}
	if req.State.Raw.IsNull() {
		return
	}

	var priorOverride types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, overridePath, &priorOverride)...)
	if resp.Diagnostics.HasError() || !priorOverride.IsNull() {
		return
	}
	reconnects, diags := utils.AppReconnectsOnUpdate(ctx, req.Plan, req.State)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() || reconnects {
		return
	}
	resp.PlanValue = req.StateValue
}

func ShowsOverride() planmodifier.String {
	return StringShowsOverridePlanModifier{}
}
