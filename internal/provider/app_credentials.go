package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/utils"
)

// keepUnchangedAppCredentials leaves auth, settings and version out of a
// lumos_app update that doesn't change them, so the API keeps the
// integration's connection instead of reconnecting it. A reconnect starts a
// sync, and returns an integration finished with browser consent to
// Partially Connected. The update is a merge patch, so the omitted fields
// keep their values, and Update copies the plan back afterwards, so state
// still records them.
//
// The generated Update calls this through
// .speakeasy/patches/internal/provider/app_resource.go.patch.
func keepUnchangedAppCredentials(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, data *AppResourceModel) {
	if data == nil {
		return
	}
	reconnects, diags := utils.AppReconnectsOnUpdate(ctx, req.Plan, req.State)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() || reconnects {
		return
	}
	data.Auth = jsontypes.NewNormalizedNull()
	data.Settings = jsontypes.NewNormalizedNull()
	data.Version = types.StringNull()
}
