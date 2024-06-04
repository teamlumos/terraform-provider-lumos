// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetActivityLogsRequest struct {
	Since  *string `queryParam:"style=form,explode=true,name=since"`
	Until  *string `queryParam:"style=form,explode=true,name=until"`
	Limit  *int64  `default:"50" queryParam:"style=form,explode=true,name=limit"`
	Offset *int64  `default:"0" queryParam:"style=form,explode=true,name=offset"`
}

func (g GetActivityLogsRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetActivityLogsRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetActivityLogsRequest) GetSince() *string {
	if o == nil {
		return nil
	}
	return o.Since
}

func (o *GetActivityLogsRequest) GetUntil() *string {
	if o == nil {
		return nil
	}
	return o.Until
}

func (o *GetActivityLogsRequest) GetLimit() *int64 {
	if o == nil {
		return nil
	}
	return o.Limit
}

func (o *GetActivityLogsRequest) GetOffset() *int64 {
	if o == nil {
		return nil
	}
	return o.Offset
}

type GetActivityLogsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	LimitOffsetPageActivityLog *shared.LimitOffsetPageActivityLog
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetActivityLogsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetActivityLogsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetActivityLogsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetActivityLogsResponse) GetLimitOffsetPageActivityLog() *shared.LimitOffsetPageActivityLog {
	if o == nil {
		return nil
	}
	return o.LimitOffsetPageActivityLog
}

func (o *GetActivityLogsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
