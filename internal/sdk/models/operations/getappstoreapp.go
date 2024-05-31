// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAppStoreAppRequest struct {
	AppID string `pathParam:"style=simple,explode=false,name=app_id"`
}

func (o *GetAppStoreAppRequest) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

type GetAppStoreAppResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	AppStoreApp *shared.AppStoreApp
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAppStoreAppResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAppStoreAppResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAppStoreAppResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAppStoreAppResponse) GetAppStoreApp() *shared.AppStoreApp {
	if o == nil {
		return nil
	}
	return o.AppStoreApp
}

func (o *GetAppStoreAppResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
