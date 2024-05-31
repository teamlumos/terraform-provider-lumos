// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type ActivityRecordsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	ActivityRecordOutput *shared.ActivityRecordOutput
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *ActivityRecordsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *ActivityRecordsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *ActivityRecordsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *ActivityRecordsResponse) GetActivityRecordOutput() *shared.ActivityRecordOutput {
	if o == nil {
		return nil
	}
	return o.ActivityRecordOutput
}

func (o *ActivityRecordsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
