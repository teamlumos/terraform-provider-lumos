// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type AddAppToAppStoreResponse struct {
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

func (o *AddAppToAppStoreResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *AddAppToAppStoreResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *AddAppToAppStoreResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *AddAppToAppStoreResponse) GetAppStoreAppSettingsOutput() *shared.AppStoreAppSettingsOutput {
	if o == nil {
		return nil
	}
	return o.AppStoreAppSettingsOutput
}

func (o *AddAppToAppStoreResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
