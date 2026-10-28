package utils

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAppReconnectsOnUpdate(t *testing.T) {
	ctx := context.Background()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"auth":         schema.StringAttribute{Optional: true, Sensitive: true, CustomType: jsontypes.NormalizedType{}},
		"settings":     schema.StringAttribute{Optional: true, Sensitive: true, CustomType: jsontypes.NormalizedType{}},
		"version":      schema.StringAttribute{Optional: true},
		"disconnected": schema.BoolAttribute{Computed: true},
	}}
	objectType := s.Type().TerraformType(ctx)
	str := func(v string) tftypes.Value { return tftypes.NewValue(tftypes.String, v) }
	none := tftypes.NewValue(tftypes.String, nil)
	unknown := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	type app struct{ auth, settings, version tftypes.Value }
	value := func(a app, disconnected bool) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{
			"auth": a.auth, "settings": a.settings, "version": a.version,
			"disconnected": tftypes.NewValue(tftypes.Bool, disconnected),
		})
	}
	connected := app{str(`{"token":"v1"}`), str(`{"region":"us"}`), str("2")}

	tests := []struct {
		name         string
		desired      app
		prior        app
		disconnected bool
		want         bool
	}{
		{name: "unchanged credentials", desired: connected, prior: connected, want: false},
		{name: "reformatted JSON", desired: app{str(`{ "token": "v1" }`), str(`{"region": "us"}`), str("2")}, prior: connected, want: false},
		{name: "changed auth", desired: app{str(`{"token":"v2"}`), connected.settings, connected.version}, prior: connected, want: true},
		{name: "changed settings", desired: app{connected.auth, str(`{"region":"eu"}`), connected.version}, prior: connected, want: true},
		{name: "removed settings", desired: app{connected.auth, none, connected.version}, prior: connected, want: true},
		{name: "changed version", desired: app{connected.auth, connected.settings, str("3")}, prior: connected, want: true},
		{name: "auth not known until apply", desired: app{unknown, connected.settings, connected.version}, prior: connected, want: true},
		{name: "disconnected app", desired: connected, prior: connected, disconnected: true, want: true},
		{name: "custom app", desired: app{none, none, none}, prior: app{none, none, none}, want: false},
		{name: "imported connected app", desired: connected, prior: app{none, none, none}, want: false},
		{name: "imported disconnected app", desired: connected, prior: app{none, none, none}, disconnected: true, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desired := tfsdk.Config{Schema: s, Raw: value(tt.desired, false)}
			prior := tfsdk.State{Schema: s, Raw: value(tt.prior, tt.disconnected)}

			got, diags := AppReconnectsOnUpdate(ctx, desired, prior)

			if diags.HasError() {
				t.Fatalf("diagnostics: %v", diags)
			}
			if got != tt.want {
				t.Errorf("AppReconnectsOnUpdate = %v, want %v", got, tt.want)
			}
		})
	}
}
