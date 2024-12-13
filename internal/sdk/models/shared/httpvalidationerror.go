// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

import (
	"errors"
	"fmt"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
)

type DetailType string

const (
	DetailTypeArrayOfValidationError DetailType = "arrayOfValidationError"
	DetailTypeStr                    DetailType = "str"
)

// Detail - The detail of the validation error
type Detail struct {
	ArrayOfValidationError []ValidationError `queryParam:"inline"`
	Str                    *string           `queryParam:"inline"`

	Type DetailType
}

func CreateDetailArrayOfValidationError(arrayOfValidationError []ValidationError) Detail {
	typ := DetailTypeArrayOfValidationError

	return Detail{
		ArrayOfValidationError: arrayOfValidationError,
		Type:                   typ,
	}
}

func CreateDetailStr(str string) Detail {
	typ := DetailTypeStr

	return Detail{
		Str:  &str,
		Type: typ,
	}
}

func (u *Detail) UnmarshalJSON(data []byte) error {

	var arrayOfValidationError []ValidationError = []ValidationError{}
	if err := utils.UnmarshalJSON(data, &arrayOfValidationError, "", true, true); err == nil {
		u.ArrayOfValidationError = arrayOfValidationError
		u.Type = DetailTypeArrayOfValidationError
		return nil
	}

	var str string = ""
	if err := utils.UnmarshalJSON(data, &str, "", true, true); err == nil {
		u.Str = &str
		u.Type = DetailTypeStr
		return nil
	}

	return fmt.Errorf("could not unmarshal `%s` into any supported union types for Detail", string(data))
}

func (u Detail) MarshalJSON() ([]byte, error) {
	if u.ArrayOfValidationError != nil {
		return utils.MarshalJSON(u.ArrayOfValidationError, "", true)
	}

	if u.Str != nil {
		return utils.MarshalJSON(u.Str, "", true)
	}

	return nil, errors.New("could not marshal union type Detail: all fields are null")
}

type HTTPValidationError struct {
	// The detail of the validation error
	Detail *Detail `json:"detail,omitempty"`
}

func (o *HTTPValidationError) GetDetail() *Detail {
	if o == nil {
		return nil
	}
	return o.Detail
}
