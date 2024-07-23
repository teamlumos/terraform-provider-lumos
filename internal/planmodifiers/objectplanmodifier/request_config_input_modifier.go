// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package objectplanmodifier

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

var _ planmodifier.Object = ObjectRequestConfigInputModifierPlanModifier{}

type ObjectRequestConfigInputModifierPlanModifier struct{}

// Description describes the plan modification in plain text formatting.
func (v ObjectRequestConfigInputModifierPlanModifier) Description(_ context.Context) string {
	return "TODO: add plan modifier description"
}

// MarkdownDescription describes the plan modification in Markdown formatting.
func (v ObjectRequestConfigInputModifierPlanModifier) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the plan modification.
func (v ObjectRequestConfigInputModifierPlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	resp.Diagnostics.AddAttributeError(
		req.Path,
		"TODO: implement plan modifier RequestConfigInputModifier logic",
		req.Path.String()+": "+v.Description(ctx),
	)
}

func RequestConfigInputModifier() planmodifier.Object {
	return ObjectRequestConfigInputModifierPlanModifier{}
}
