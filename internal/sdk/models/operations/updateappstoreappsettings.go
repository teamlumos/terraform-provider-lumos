// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type UpdateAppStoreAppSettingsRequest struct {
	AppID                    string                          `pathParam:"style=simple,explode=false,name=app_id"`
	AppStoreAppSettingsInput shared.AppStoreAppSettingsInput `request:"mediaType=application/json"`
}

func (o *UpdateAppStoreAppSettingsRequest) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *UpdateAppStoreAppSettingsRequest) GetAppStoreAppSettingsInput() shared.AppStoreAppSettingsInput {
	if o == nil {
		return shared.AppStoreAppSettingsInput{}
	}
	return o.AppStoreAppSettingsInput
}

type UpdateAppStoreAppSettingsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	AppStoreAppSettingsOutput *shared.AppStoreAppSettingsOutput
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *UpdateAppStoreAppSettingsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdateAppStoreAppSettingsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdateAppStoreAppSettingsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdateAppStoreAppSettingsResponse) GetAppStoreAppSettingsOutput() *shared.AppStoreAppSettingsOutput {
	if o == nil {
		return nil
	}
	return o.AppStoreAppSettingsOutput
}

func (o *UpdateAppStoreAppSettingsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
