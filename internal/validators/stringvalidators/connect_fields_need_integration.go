package stringvalidators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

var _ validator.String = StringConnectFieldsNeedIntegrationValidator{}

// connectFields connect an integration, so a custom app can't set them.
var connectFields = []string{"auth", "settings", "version"}

// StringConnectFieldsNeedIntegrationValidator fails a lumos_app config that
// sets auth, settings or version without an integration to connect: with
// app_class_id omitted, or set to "custom", the value a custom app reads
// back with.
type StringConnectFieldsNeedIntegrationValidator struct{}

// Description describes the validation in plain text formatting.
func (v StringConnectFieldsNeedIntegrationValidator) Description(_ context.Context) string {
	return `a custom app (app_class_id omitted or "custom") can't set auth, settings or version`
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v StringConnectFieldsNeedIntegrationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v StringConnectFieldsNeedIntegrationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() {
		return
	}
	if !req.ConfigValue.IsNull() && req.ConfigValue.ValueString() != utils.CustomAppClassID {
		return
	}
	for _, name := range connectFields {
		var value attr.Value
		diags := req.Config.GetAttribute(ctx, path.Root(name), &value)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		if !value.IsNull() {
			resp.Diagnostics.AddAttributeError(path.Root(name), "Field needs an integration",
				fmt.Sprintf("%s only applies when connecting an integration, and this is a custom app. "+
					"Remove %s, or set app_class_id to the integration to connect.", name, name))
		}
	}
}

func ConnectFieldsNeedIntegration() validator.String {
	return StringConnectFieldsNeedIntegrationValidator{}
}
