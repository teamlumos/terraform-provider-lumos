// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetRequest struct {
	// Filters preapproval rules by the ID of the app to which they belong.
	AppID *string `queryParam:"style=form,explode=true,name=app_id"`
	Page  *int64  `default:"1" queryParam:"style=form,explode=true,name=page"`
	Size  *int64  `default:"50" queryParam:"style=form,explode=true,name=size"`
}

func (g GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetRequest) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

type GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	PagePreApprovalRuleOutput *shared.PagePreApprovalRuleOutput
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetResponse) GetPagePreApprovalRuleOutput() *shared.PagePreApprovalRuleOutput {
	if o == nil {
		return nil
	}
	return o.PagePreApprovalRuleOutput
}

func (o *GetAppstorePreApprovalRulesForAppAppstorePreApprovalRulesGetResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
