package stringvalidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = StringLeaveOutToClearValidator{}

// StringLeaveOutToClearValidator fails the plan when an attribute is set to
// "". The API doesn't keep "" for these attributes: it either stores null,
// so every plan would show a change back to "", or refuses the value.
// Leaving the attribute out clears it instead.
type StringLeaveOutToClearValidator struct{}

// Description describes the validation in plain text formatting.
func (v StringLeaveOutToClearValidator) Description(_ context.Context) string {
	return `value must not be ""; leave the attribute out to clear it`
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v StringLeaveOutToClearValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v StringLeaveOutToClearValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() || req.ConfigValue.ValueString() != "" {
		return
	}
	// Only the attribute's own name: inside a set, the full path names the
	// element by its whole value. The diagnostic already points at the block.
	name := req.Path.String()
	if step, _ := req.Path.Steps().LastStep(); step != nil {
		name = step.String()
	}
	resp.Diagnostics.AddAttributeError(req.Path, fmt.Sprintf(`%s can't be ""`, name),
		fmt.Sprintf(`Lumos doesn't keep "" as a value. Leave %s out to clear it.`, name))
}

func LeaveOutToClear() validator.String {
	return StringLeaveOutToClearValidator{}
}
