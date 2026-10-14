package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// AppReconnectsOnUpdate reports whether updating a lumos_app reconnects its
// integration. The update sends auth and settings whenever the config sets
// them, and the API reconnects and syncs whenever it receives credentials,
// so every update of a connected app reconnects and a custom app's never
// does.
func AppReconnectsOnUpdate(ctx context.Context, config tfsdk.Config) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	for _, name := range []string{"auth", "settings"} {
		var value jsontypes.Normalized
		diags.Append(config.GetAttribute(ctx, path.Root(name), &value)...)
		if !value.IsNull() {
			return true, diags
		}
	}
	return false, diags
}
