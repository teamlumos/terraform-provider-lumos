package stringplanmodifier

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.String = StringRequireOnceSetPlanModifier{}

// StringRequireOnceSetPlanModifier fails the plan when the config leaves out
// a field that Lumos holds a value for and can't clear, such as a workflow's
// app_id once the workflow is created with an app or bound to one. A config
// that leaves the field out means "cleared", which Lumos can't do, so the
// config has to state the value instead. The state is read from Lumos when
// the plan refreshes, so a value Lumos assigned is caught too.
type StringRequireOnceSetPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringRequireOnceSetPlanModifier) Description(_ context.Context) string {
	return "Fails the plan when the config leaves out a value Lumos holds and can't clear."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringRequireOnceSetPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyString performs the plan modification.
func (v StringRequireOnceSetPlanModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || !req.ConfigValue.IsNull() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	resp.Diagnostics.AddAttributeError(req.Path, fmt.Sprintf("%s can't be cleared", req.Path),
		fmt.Sprintf("Lumos holds %s = %q and can't clear or change it once it's set. Add %s = %q to the config.",
			req.Path, req.StateValue.ValueString(), req.Path, req.StateValue.ValueString()))
}

func RequireOnceSet() planmodifier.String {
	return StringRequireOnceSetPlanModifier{}
}
