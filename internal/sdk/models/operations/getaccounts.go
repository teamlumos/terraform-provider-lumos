// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAccountsRequest struct {
	AppID *string `queryParam:"style=form,explode=true,name=app_id"`
	Page  *int64  `default:"1" queryParam:"style=form,explode=true,name=page"`
	Size  *int64  `default:"50" queryParam:"style=form,explode=true,name=size"`
}

func (g GetAccountsRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAccountsRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAccountsRequest) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *GetAccountsRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *GetAccountsRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

type GetAccountsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	PageAccount *shared.PageAccount
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAccountsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAccountsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAccountsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAccountsResponse) GetPageAccount() *shared.PageAccount {
	if o == nil {
		return nil
	}
	return o.PageAccount
}

func (o *GetAccountsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
