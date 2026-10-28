package utils

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// attributeGetter is a tfsdk.Config, tfsdk.Plan or tfsdk.State.
type attributeGetter interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

// AppReconnectsOnUpdate reports whether updating a lumos_app reconnects its
// integration, comparing the desired values, from the config or plan, with
// the prior state. The API reconnects and syncs whenever it receives
// credentials, so the update sends auth, settings and version only when
// this is true, and sends all three then, since the API needs auth with
// either of the others.
//
// An update reconnects when the desired auth, settings or version differs
// from state or isn't known yet, or when the app is disconnected. A custom
// app, which sets none of them, never reconnects. Neither does a connected
// app with no auth in state, which was imported: its first update records
// the configured credentials without sending them, so it keeps a connection
// finished with browser consent.
func AppReconnectsOnUpdate(ctx context.Context, desired, prior attributeGetter) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	var auth, settings, priorAuth, priorSettings jsontypes.Normalized
	var version, priorVersion types.String
	var disconnected types.Bool
	diags.Append(desired.GetAttribute(ctx, path.Root("auth"), &auth)...)
	diags.Append(desired.GetAttribute(ctx, path.Root("settings"), &settings)...)
	diags.Append(desired.GetAttribute(ctx, path.Root("version"), &version)...)
	diags.Append(prior.GetAttribute(ctx, path.Root("auth"), &priorAuth)...)
	diags.Append(prior.GetAttribute(ctx, path.Root("settings"), &priorSettings)...)
	diags.Append(prior.GetAttribute(ctx, path.Root("version"), &priorVersion)...)
	diags.Append(prior.GetAttribute(ctx, path.Root("disconnected"), &disconnected)...)
	if diags.HasError() {
		return false, diags
	}

	if auth.IsNull() && settings.IsNull() && version.IsNull() {
		return false, diags
	}
	if auth.IsUnknown() || settings.IsUnknown() || version.IsUnknown() || disconnected.ValueBool() {
		return true, diags
	}
	if priorAuth.IsNull() {
		return false, diags
	}
	for _, pair := range [][2]jsontypes.Normalized{{auth, priorAuth}, {settings, priorSettings}} {
		equal, equalDiags := jsonEqual(ctx, pair[0], pair[1])
		diags.Append(equalDiags...)
		if equalDiags.HasError() || !equal {
			return true, diags
		}
	}
	return !version.Equal(priorVersion), diags
}

// jsonEqual compares two JSON values, ignoring whitespace and key order.
func jsonEqual(ctx context.Context, a, b jsontypes.Normalized) (bool, diag.Diagnostics) {
	if a.IsNull() || b.IsNull() {
		return a.IsNull() && b.IsNull(), nil
	}
	return a.StringSemanticEquals(ctx, b)
}
