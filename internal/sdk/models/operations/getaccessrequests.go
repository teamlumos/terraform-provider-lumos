// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAccessRequestsRequest struct {
	// Filters requests by the ID of the target user.
	TargetUserID *string `queryParam:"style=form,explode=true,name=target_user_id"`
	// Filters requests by the ID of the requesting user.
	RequesterUserID *string `queryParam:"style=form,explode=true,name=requester_user_id"`
	// Filters requests by the ID of the user.
	UserID *string `queryParam:"style=form,explode=true,name=user_id"`
	// Filters requests by their status. Possible values: 'PENDING', 'PENDING_MANAGER_APPROVAL', 'MANAGER_APPROVED', 'MANAGER_DENIED', 'PENDING_APPROVAL', 'APPROVED', 'DENIED', 'EXPIRED', 'CANCELLED', 'PENDING_PROVISIONING', 'PENDING_MANUAL_PROVISIONING', 'DENIED_PROVISIONING', 'PROVISIONED', 'PENDING_DEPROVISIONING', 'PENDING_MANUAL_DEPROVISIONING', 'TIME_BASED_EXPIRED', 'COMPLETED', 'REVERTING', 'REVERTED'
	Statuses []shared.SupportRequestStatus `queryParam:"style=form,explode=true,name=statuses"`
	// Sort access requests ascending (ASC) or descending (DESC) by created date.
	Sort *string `default:"ASC" queryParam:"style=form,explode=true,name=sort"`
	// Page number
	Page *int64 `default:"1" queryParam:"style=form,explode=true,name=page"`
	// Page size
	Size *int64 `default:"50" queryParam:"style=form,explode=true,name=size"`
}

func (g GetAccessRequestsRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAccessRequestsRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAccessRequestsRequest) GetTargetUserID() *string {
	if o == nil {
		return nil
	}
	return o.TargetUserID
}

func (o *GetAccessRequestsRequest) GetRequesterUserID() *string {
	if o == nil {
		return nil
	}
	return o.RequesterUserID
}

func (o *GetAccessRequestsRequest) GetUserID() *string {
	if o == nil {
		return nil
	}
	return o.UserID
}

func (o *GetAccessRequestsRequest) GetStatuses() []shared.SupportRequestStatus {
	if o == nil {
		return nil
	}
	return o.Statuses
}

func (o *GetAccessRequestsRequest) GetSort() *string {
	if o == nil {
		return nil
	}
	return o.Sort
}

func (o *GetAccessRequestsRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *GetAccessRequestsRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

type GetAccessRequestsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	PageAccessRequest *shared.PageAccessRequest
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAccessRequestsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAccessRequestsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAccessRequestsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAccessRequestsResponse) GetPageAccessRequest() *shared.PageAccessRequest {
	if o == nil {
		return nil
	}
	return o.PageAccessRequest
}

func (o *GetAccessRequestsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
