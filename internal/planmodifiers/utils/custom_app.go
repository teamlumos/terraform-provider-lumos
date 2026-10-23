package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomAppClassID is the app_class_id the API reads back for a custom app.
const CustomAppClassID = "custom"

// AppConfigIsCustom reports whether a lumos_app config is for a custom app:
// one that doesn't connect an integration through app_class_id.
func AppConfigIsCustom(ctx context.Context, config tfsdk.Config) (bool, diag.Diagnostics) {
	var appClassID types.String
	diags := config.GetAttribute(ctx, path.Root("app_class_id"), &appClassID)
	return appClassID.IsNull(), diags
}
