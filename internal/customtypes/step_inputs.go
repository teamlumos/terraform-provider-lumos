package customtypes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StepInputsType is jsontypes.NormalizedType, except that an object key set
// to null counts as the same as a missing key.
type StepInputsType struct{ jsontypes.NormalizedType }

var _ basetypes.StringTypable = StepInputsType{}

// The plugin framework requires these on every custom type. The embedded
// NormalizedType would otherwise answer them with its own type and values.
func (t StepInputsType) String() string                       { return "StepInputsType" }
func (t StepInputsType) ValueType(context.Context) attr.Value { return StepInputs{} }
func (t StepInputsType) Equal(o attr.Type) bool               { _, ok := o.(StepInputsType); return ok }

func (t StepInputsType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return StepInputs{Normalized: jsontypes.Normalized{StringValue: in}}, nil
}

func (t StepInputsType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	v, err := t.NormalizedType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	n, ok := v.(jsontypes.Normalized)
	if !ok {
		return nil, fmt.Errorf("unexpected value type %T", v)
	}
	return StepInputs{Normalized: n}, nil
}

// StepInputs is a step's inputs: a JSON string that keeps the Normalized
// value's JSON validation and differs only in how it compares.
type StepInputs struct{ jsontypes.Normalized }

var _ basetypes.StringValuableWithSemanticEquals = StepInputs{}

func (v StepInputs) Type(context.Context) attr.Type { return StepInputsType{} }

// Equal is overridden because Normalized's Equal only accepts a Normalized.
func (v StepInputs) Equal(o attr.Value) bool {
	other, ok := o.(StepInputs)
	return ok && v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals tells Terraform that the config's inputs and the
// inputs the API read back are the same value when they differ only in key
// order, spacing, or keys set to null.
func (v StepInputs) StringSemanticEquals(_ context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	nv, ok := newValuable.(StepInputs)
	if !ok {
		diags.AddError("Semantic Equality Check Error", fmt.Sprintf("unexpected value type %T", newValuable))
		return false, diags
	}
	a, err := withoutNulls(v.ValueString())
	if err != nil {
		diags.AddError("Semantic Equality Check Error", err.Error())
		return false, diags
	}
	b, err := withoutNulls(nv.ValueString())
	if err != nil {
		diags.AddError("Semantic Equality Check Error", err.Error())
		return false, diags
	}
	return a == b, diags
}

// withoutNulls re-encodes the JSON with null-valued object keys removed at
// every depth. Numbers stay as written (UseNumber); key order is sorted by
// encoding/json.
func withoutNulls(s string) (string, error) {
	dec := json.NewDecoder(bytes.NewReader([]byte(s)))
	dec.UseNumber()
	var x any
	if err := dec.Decode(&x); err != nil {
		return "", err
	}
	out, err := json.Marshal(strip(x))
	return string(out), err
}

func strip(x any) any {
	switch t := x.(type) {
	case map[string]any:
		m := make(map[string]any, len(t))
		for k, v := range t {
			if v != nil {
				m[k] = strip(v)
			}
		}
		return m
	case []any:
		for i := range t {
			t[i] = strip(t[i])
		}
		return t
	}
	return x
}
