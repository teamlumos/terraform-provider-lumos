package stringplanmodifier

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUnknownOnDescriptionOverrideChange(t *testing.T) {
	ctx := context.Background()
	s := schema.Schema{Attributes: map[string]schema.Attribute{
		"description_override": schema.StringAttribute{Optional: true},
		"description":          schema.StringAttribute{Computed: true},
	}}
	objectType := s.Type().TerraformType(ctx)
	entitlement := func(override, description tftypes.Value) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{"description_override": override, "description": description})
	}
	str := func(s string) tftypes.Value { return tftypes.NewValue(tftypes.String, s) }
	none := tftypes.NewValue(tftypes.String, nil)
	unknown := tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	synced := str("All engineers")

	tests := []struct {
		name      string
		config    tftypes.Value
		state     tftypes.Value
		planValue types.String
		want      types.String
	}{
		{name: "no override either side keeps the integration's description", config: entitlement(none, none), state: entitlement(none, synced), planValue: types.StringUnknown(), want: types.StringValue("All engineers")},
		{name: "an unchanged override keeps the description", config: entitlement(str("Ask first"), none), state: entitlement(str("Ask first"), str("Ask first")), planValue: types.StringUnknown(), want: types.StringValue("Ask first")},
		{name: "setting an override stays unknown", config: entitlement(str("Ask first"), none), state: entitlement(none, synced), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "changing an override stays unknown", config: entitlement(str("Ask IT"), none), state: entitlement(str("Ask first"), str("Ask first")), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "clearing an override stays unknown", config: entitlement(none, none), state: entitlement(str("Ask first"), str("Ask first")), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "an override not known until apply stays unknown", config: entitlement(unknown, none), state: entitlement(none, synced), planValue: types.StringUnknown(), want: types.StringUnknown()},
		{name: "create stays unknown", config: entitlement(none, none), state: tftypes.NewValue(objectType, nil), planValue: types.StringUnknown(), want: types.StringUnknown()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateValue := types.StringNull()
			if !tt.state.IsNull() {
				var description types.String
				state := tfsdk.State{Schema: s, Raw: tt.state}
				state.GetAttribute(ctx, path.Root("description"), &description)
				stateValue = description
			}
			req := planmodifier.StringRequest{
				Config:     tfsdk.Config{Schema: s, Raw: tt.config},
				Plan:       tfsdk.Plan{Schema: s, Raw: entitlement(none, unknown)},
				State:      tfsdk.State{Schema: s, Raw: tt.state},
				PlanValue:  tt.planValue,
				StateValue: stateValue,
			}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

			UnknownOnDescriptionOverrideChange().PlanModifyString(ctx, req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(tt.want) {
				t.Errorf("plan value = %v, want %v", resp.PlanValue, tt.want)
			}
		})
	}
}
