package stringplanmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUnknownOnReconnect(t *testing.T) {
	ctx := context.Background()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"auth":        schema.StringAttribute{Optional: true, Sensitive: true, CustomType: jsontypes.NormalizedType{}},
		"settings":    schema.StringAttribute{Optional: true, Sensitive: true, CustomType: jsontypes.NormalizedType{}},
		"sync_status": schema.StringAttribute{Computed: true},
	}}
	objectType := s.Type().TerraformType(ctx)
	app := func(auth, settings, syncStatus tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{"auth": auth, "settings": settings, "sync_status": syncStatus})
	}
	none := tftypes.NewValue(tftypes.String, nil)
	unknown := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	creds := tftypes.NewValue(tftypes.String, `{"api_key":"x"}`)
	success := tftypes.NewValue(tftypes.String, "success")

	tests := []struct {
		name      string
		config    tftypes.Value
		state     tftypes.Value
		planValue types.String
		want      types.String
	}{
		{name: "connected app update stays unknown", config: app(creds, none, none), state: app(creds, none, success), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "settings alone count as credentials", config: app(none, creds, none), state: app(none, creds, success), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "custom app update keeps the current value", config: app(none, none, none), state: app(none, none, success), planValue: types.StringUnknown(), want: types.StringValue("success")},
		{name: "create stays unknown", config: app(none, none, none), state: tftypes.NewValue(objectType, nil), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "a known plan value is left alone", config: app(none, none, none), state: app(none, none, success), planValue: types.StringValue("failed"), want: types.StringValue("failed")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateValue := types.StringNull()
			if !tt.state.IsNull() {
				stateValue = types.StringValue("success")
			}
			req := planmodifier.StringRequest{
				Config:     tfsdk.Config{Schema: s, Raw: tt.config},
				Plan:       tfsdk.Plan{Schema: s, Raw: app(none, none, unknown)},
				State:      tfsdk.State{Schema: s, Raw: tt.state},
				PlanValue:  tt.planValue,
				StateValue: stateValue,
			}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

			UnknownOnReconnect().PlanModifyString(ctx, req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(tt.want) {
				t.Errorf("plan value = %v, want %v", resp.PlanValue, tt.want)
			}
		})
	}
}
