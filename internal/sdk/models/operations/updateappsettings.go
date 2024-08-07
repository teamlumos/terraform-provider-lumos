// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type UpdateAppSettingsRequest struct {
	AppID           string                 `pathParam:"style=simple,explode=false,name=app_id"`
	AppSettingInput shared.AppSettingInput `request:"mediaType=application/json"`
}

func (o *UpdateAppSettingsRequest) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *UpdateAppSettingsRequest) GetAppSettingInput() shared.AppSettingInput {
	if o == nil {
		return shared.AppSettingInput{}
	}
	return o.AppSettingInput
}

type UpdateAppSettingsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	AppSettingOutput *shared.AppSettingOutput
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *UpdateAppSettingsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateAppSettingsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateAppSettingsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateAppSettingsResponse) GetAppSettingOutput() *shared.AppSettingOutput {
	if o == nil {
		return nil
	}
	return o.AppSettingOutput
}

func (o *UpdateAppSettingsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
