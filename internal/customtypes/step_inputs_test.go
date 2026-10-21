package customtypes

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
)

func TestNullKeys(t *testing.T) {
	cases := []struct {
		name       string
		config     string
		api        string
		jsontypes  bool
		stepInputs bool
	}{
		{"top-level null vs missing", `{"fallback_workflow_id":null,"branches":[]}`, `{"branches":[]}`, false, true},
		{"nested null vs missing", `{"branches":[{"workflow_id":"w","condition":{"attribute":"a","equals":null}}]}`, `{"branches":[{"workflow_id":"w","condition":{"attribute":"a"}}]}`, false, true},
		{"key order and spacing", `{"b":1, "a":2}`, `{"a":2,"b":1}`, true, true},
		{"real change still differs", `{"stage_number":"1"}`, `{"stage_number":1}`, false, false},
		{"real value vs null still differs", `{"x":"y"}`, `{"x":null}`, false, false},
		{"array null element kept", `{"x":[null]}`, `{"x":[]}`, false, false},
	}
	ctx := context.Background()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			j, _ := jsontypes.NewNormalizedValue(c.config).StringSemanticEquals(ctx, jsontypes.NewNormalizedValue(c.api))
			if j != c.jsontypes {
				t.Errorf("jsontypes: got %v want %v", j, c.jsontypes)
			}
			cfg := StepInputs{Normalized: jsontypes.NewNormalizedValue(c.config)}
			api := StepInputs{Normalized: jsontypes.NewNormalizedValue(c.api)}
			s, d := cfg.StringSemanticEquals(ctx, api)
			if d.HasError() {
				t.Fatal(d)
			}
			if s != c.stepInputs {
				t.Errorf("StepInputs: got %v want %v", s, c.stepInputs)
			}
		})
	}
}
