// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAccessRequestRequest struct {
	AccessRequestID string `pathParam:"style=simple,explode=false,name=access_request_id"`
	Page            *int64 `default:"1" queryParam:"style=form,explode=true,name=page"`
	Size            *int64 `default:"50" queryParam:"style=form,explode=true,name=size"`
}

func (g GetAccessRequestRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAccessRequestRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAccessRequestRequest) GetAccessRequestID() string {
	if o == nil {
		return ""
	}
	return o.AccessRequestID
}

func (o *GetAccessRequestRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *GetAccessRequestRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

type GetAccessRequestResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	AccessRequest *shared.AccessRequest
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAccessRequestResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAccessRequestResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAccessRequestResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAccessRequestResponse) GetAccessRequest() *shared.AccessRequest {
	if o == nil {
		return nil
	}
	return o.AccessRequest
}

func (o *GetAccessRequestResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
