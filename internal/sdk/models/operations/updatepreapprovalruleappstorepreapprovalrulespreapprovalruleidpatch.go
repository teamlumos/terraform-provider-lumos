// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchRequest struct {
	ID                         string                            `pathParam:"style=simple,explode=false,name=pre_approval_rule_id"`
	PreApprovalRuleUpdateInput shared.PreApprovalRuleUpdateInput `request:"mediaType=application/json"`
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchRequest) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchRequest) GetPreApprovalRuleUpdateInput() shared.PreApprovalRuleUpdateInput {
	if o == nil {
		return shared.PreApprovalRuleUpdateInput{}
	}
	return o.PreApprovalRuleUpdateInput
}

type UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	PreApprovalRuleOutput *shared.PreApprovalRuleOutput
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchResponse) GetPreApprovalRuleOutput() *shared.PreApprovalRuleOutput {
	if o == nil {
		return nil
	}
	return o.PreApprovalRuleOutput
}

func (o *UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
