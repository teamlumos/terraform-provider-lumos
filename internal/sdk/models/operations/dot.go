// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type DotRequest struct {
	AppID string `pathParam:"style=simple,explode=false,name=app_id"`
}

func (o *DotRequest) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

type DotResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	Dot *string
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *DotResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DotResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DotResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DotResponse) GetDot() *string {
	if o == nil {
		return nil
	}
	return o.Dot
}

func (o *DotResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
