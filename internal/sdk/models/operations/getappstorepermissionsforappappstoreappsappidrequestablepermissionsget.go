// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest struct {
	// Filters requestable permissions by the ID of the app to which they belong.
	AppID string `pathParam:"style=simple,explode=false,name=app_id"`
	//     Searches permissions by the permission's group name,
	//     request configuration name, or specific integration ID.
	//
	SearchTerm *string `queryParam:"style=form,explode=true,name=search_term"`
	// Search filter should be an exact match.
	ExactMatch *bool `default:"false" queryParam:"style=form,explode=true,name=exact_match"`
	//     Filters permissions by visibility in the AppStore.
	//
	InAppStore *bool `queryParam:"style=form,explode=true,name=in_app_store"`
	// Include inherited configurations from parent app.
	IncludeInheritedConfigs *bool `default:"true" queryParam:"style=form,explode=true,name=include_inherited_configs"`
	// Page number
	Page *int64 `default:"1" queryParam:"style=form,explode=true,name=page"`
	// Page size
	Size *int64 `default:"50" queryParam:"style=form,explode=true,name=size"`
}

func (g GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetAppID() string {
	if o == nil {
		return ""
	}
	return o.AppID
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetSearchTerm() *string {
	if o == nil {
		return nil
	}
	return o.SearchTerm
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetExactMatch() *bool {
	if o == nil {
		return nil
	}
	return o.ExactMatch
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetInAppStore() *bool {
	if o == nil {
		return nil
	}
	return o.InAppStore
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetIncludeInheritedConfigs() *bool {
	if o == nil {
		return nil
	}
	return o.IncludeInheritedConfigs
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

type GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	PageRequestablePermissionOutput *shared.PageRequestablePermissionOutput
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetResponse) GetPageRequestablePermissionOutput() *shared.PageRequestablePermissionOutput {
	if o == nil {
		return nil
	}
	return o.PageRequestablePermissionOutput
}

func (o *GetAppstorePermissionsForAppAppstoreAppsAppIDRequestablePermissionsGetResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
