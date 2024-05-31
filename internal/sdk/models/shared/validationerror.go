// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"errors"
	"fmt"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/internal/utils"
)

type LocType string

const (
	LocTypeStr     LocType = "str"
	LocTypeInteger LocType = "integer"
)

type Loc struct {
	Str     *string
	Integer *int64

	Type LocType
}

func CreateLocStr(str string) Loc {
	typ := LocTypeStr

	return Loc{
		Str:  &str,
		Type: typ,
	}
}

func CreateLocInteger(integer int64) Loc {
	typ := LocTypeInteger

	return Loc{
		Integer: &integer,
		Type:    typ,
	}
}

func (u *Loc) UnmarshalJSON(data []byte) error {

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = LocTypeStr
		return nil
	}

	var integer int64 = int64(0)
	if err := utils.UnmarshalJSON(data, &integer, "", true, true); err == nil {
		u.Integer = &integer
		u.Type = LocTypeInteger
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for Loc", string(data))
}

func (u Loc) MarshalJSON() ([]byte, error) {
	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	if u.Integer != nil {
		return utils.MarshalJSON(u.Integer, "", true)
	}

	return nil, errors.New("could not marshal union type Loc: all fields are null")
}

type ValidationError struct {
	Loc  []Loc  `json:"loc"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

func (o *ValidationError) GetLoc() []Loc {
	if o == nil {
		return []Loc{}
	}
	return o.Loc
}

func (o *ValidationError) GetMsg() string {
	if o == nil {
		return ""
	}
	return o.Msg
}

func (o *ValidationError) GetType() string {
	if o == nil {
		return ""
	}
	return o.Type
}
