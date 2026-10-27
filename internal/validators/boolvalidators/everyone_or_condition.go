package boolvalidators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Bool = BoolEveryoneOrConditionValidator{}

// BoolEveryoneOrConditionValidator fails the plan unless exactly one of
// is_everyone_condition = true and access_condition is set. Lumos refuses
// both and neither with a 400. is_everyone_condition = false next to an
// access_condition is valid, so ConflictsWith doesn't fit. It runs on
// is_everyone_condition but reads access_condition from the config, so it
// also fails when both are left out.
type BoolEveryoneOrConditionValidator struct{}

// Description describes the validation in plain text formatting.
func (v BoolEveryoneOrConditionValidator) Description(_ context.Context) string {
	return "exactly one of access_condition and is_everyone_condition = true must be set"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v BoolEveryoneOrConditionValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v BoolEveryoneOrConditionValidator) ValidateBool(ctx context.Context, req validator.BoolRequest, resp *validator.BoolResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	var condition jsontypes.Normalized
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("access_condition"), &condition)...)
	if resp.Diagnostics.HasError() || condition.IsUnknown() {
		return
	}
	everyone := req.ConfigValue.ValueBool()
	switch {
	case everyone && !condition.IsNull():
		resp.Diagnostics.AddAttributeError(req.Path, "access_condition and is_everyone_condition = true can't both be set",
			"A policy applies to everyone or to the identities matching access_condition, not both. Use access_condition or is_everyone_condition = true.")
	case !everyone && condition.IsNull():
		resp.Diagnostics.AddAttributeError(req.Path, "access_condition or is_everyone_condition = true is required",
			"Set access_condition to the condition identities must match, or set is_everyone_condition = true to apply the policy to everyone.")
	}
}

func EveryoneOrCondition() validator.Bool {
	return BoolEveryoneOrConditionValidator{}
}
