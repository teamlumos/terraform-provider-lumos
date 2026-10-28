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
		"app_class_id": schema.StringAttribute{Optional: true},
		"auth":         schema.StringAttribute{Optional: true, Sensitive: true, CustomType: jsontypes.NormalizedType{}},
		"settings":     schema.StringAttribute{Optional: true, Sensitive: true, CustomType: jsontypes.NormalizedType{}},
		"version":      schema.StringAttribute{Optional: true},
		"disconnected": schema.BoolAttribute{Computed: true},
		"sync_status":  schema.StringAttribute{Computed: true},
	}}
	objectType := s.Type().TerraformType(ctx)
	none := tftypes.NewValue(tftypes.String, nil)
	creds := tftypes.NewValue(tftypes.String, `{"api_key":"x"}`)
	rotated := tftypes.NewValue(tftypes.String, `{"api_key":"y"}`)
	slack := tftypes.NewValue(tftypes.String, "slack.com")
	custom := tftypes.NewValue(tftypes.String, "custom")
	type app struct {
		classID, auth, settings tftypes.Value
		disconnected            bool
	}
	value := func(a app, syncStatus tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{
			"app_class_id": a.classID, "auth": a.auth, "settings": a.settings,
			"version":      tftypes.NewValue(tftypes.String, nil),
			"disconnected": tftypes.NewValue(tftypes.Bool, a.disconnected),
			"sync_status":  syncStatus,
		})
	}
	connected := app{slack, creds, creds, false}

	tests := []struct {
		name      string
		config    app
		state     *app
		planValue types.String
		want      types.String
	}{
		{name: "credential change stays unknown", config: app{slack, rotated, creds, false}, state: &connected, planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "settings change stays unknown", config: app{slack, creds, rotated, false}, state: &connected, planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "update without a credential change keeps the current value", config: connected, state: &connected, planValue: types.StringUnknown(), want: types.StringValue("success")},
		{name: "custom app update keeps the current value", config: app{custom, none, none, false}, state: &app{custom, none, none, false}, planValue: types.StringUnknown(), want: types.StringValue("success")},
		{name: "create stays unknown", config: connected, planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "a known plan value is left alone", config: connected, state: &connected, planValue: types.StringValue("success"), want: types.StringValue("success")},
		{name: "a disconnected integration plans a reconnect", config: connected, state: &app{slack, creds, creds, true}, planValue: types.StringValue("success"), want: types.StringUnknown()},
		{name: "a disconnected custom app is left alone", config: app{custom, none, none, false}, state: &app{custom, none, none, true}, planValue: types.StringValue("success"), want: types.StringValue("success")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := tftypes.NewValue(objectType, nil)
			stateValue := types.StringNull()
			if tt.state != nil {
				state = value(*tt.state, tftypes.NewValue(tftypes.String, "success"))
				stateValue = types.StringValue("success")
			}
			config := value(tt.config, none)
			plan := value(tt.config, tftypes.NewValue(tftypes.String, tftypes.UnknownValue))
			req := planmodifier.StringRequest{
				Config:     tfsdk.Config{Schema: s, Raw: config},
				Plan:       tfsdk.Plan{Schema: s, Raw: plan},
				State:      tfsdk.State{Schema: s, Raw: state},
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
