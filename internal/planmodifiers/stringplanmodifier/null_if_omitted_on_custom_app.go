package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

var _ planmodifier.String = StringNullIfOmittedOnCustomAppPlanModifier{}

// StringNullIfOmittedOnCustomAppPlanModifier plans null for a lumos_app
// field that a custom app's config leaves out, so removing it clears it and
// the plan shows that. A connected app can't set the field, and shows its
// integration's value, so its plan keeps the current value.
type StringNullIfOmittedOnCustomAppPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringNullIfOmittedOnCustomAppPlanModifier) Description(_ context.Context) string {
	return "Plans null when a custom app's config leaves the field out."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringNullIfOmittedOnCustomAppPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyString performs the plan modification.
func (v StringNullIfOmittedOnCustomAppPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	custom, diags := utils.AppConfigIsCustom(ctx, req.Config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	if custom {
		resp.PlanValue = types.StringNull()
	} else if !req.State.Raw.IsNull() {
		resp.PlanValue = req.StateValue
	}
}

func NullIfOmittedOnCustomApp() planmodifier.String {
	return StringNullIfOmittedOnCustomAppPlanModifier{}
}
