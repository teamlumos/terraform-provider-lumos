package stringplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = StringUnknownOnDescriptionOverrideChangePlanModifier{}

// StringUnknownOnDescriptionOverrideChangePlanModifier plans a lumos_entitlement's
// description, which the API reads back as description_override if set and
// otherwise the integration's description. It stays unknown when the update
// changes description_override, and keeps its current value otherwise.
type StringUnknownOnDescriptionOverrideChangePlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v StringUnknownOnDescriptionOverrideChangePlanModifier) Description(_ context.Context) string {
	return "Keeps the current value unless the update changes description_override."
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v StringUnknownOnDescriptionOverrideChangePlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// PlanModifyString performs the plan modification.
func (v StringUnknownOnDescriptionOverrideChangePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() || !req.PlanValue.IsUnknown() {
		return
	}
	var configured, current types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description_override"), &configured)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("description_override"), &current)...)
	if resp.Diagnostics.HasError() || configured.IsUnknown() || !configured.Equal(current) {
		return
	}
	resp.PlanValue = req.StateValue
}

func UnknownOnDescriptionOverrideChange() planmodifier.String {
	return StringUnknownOnDescriptionOverrideChangePlanModifier{}
}
