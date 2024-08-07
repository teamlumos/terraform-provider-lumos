// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteRequest struct {
	ID string `pathParam:"style=simple,explode=false,name=pre_approval_rule_id"`
}

func (o *DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

type DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
