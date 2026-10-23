package stringplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

var _ planmodifier.String = StringRejectIntegrationRemovalPlanModifier{}

// StringRejectIntegrationRemovalPlanModifier fails the plan when a connected
// app's config drops app_class_id. There's no empty value to clear it to: a
// lumos_app without app_class_id is a custom app, so the change would mean
// replacing the integration with a new custom app.
type StringRejectIntegrationRemovalPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringRejectIntegrationRemovalPlanModifier) Description(_ context.Context) string {
	return "Fails the plan when a connected app's config removes app_class_id."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringRejectIntegrationRemovalPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyString performs the plan modification.
func (v StringRejectIntegrationRemovalPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.ValueString() == utils.CustomAppClassID {
		return
	}
	resp.Diagnostics.AddAttributeError(req.Path, "app_class_id can't be removed",
		fmt.Sprintf("This app is connected to the %q integration, and removing app_class_id would make it a custom app. "+
			"Put app_class_id back, or to stop managing the app, remove the resource.", req.StateValue.ValueString()))
}

func RejectIntegrationRemoval() planmodifier.String {
	return StringRejectIntegrationRemovalPlanModifier{}
}
