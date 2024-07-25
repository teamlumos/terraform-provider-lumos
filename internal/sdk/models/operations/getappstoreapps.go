// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAppStoreAppsRequest struct {
	// Filters apps by the ID of the app.
	AppID *string `queryParam:"style=form,explode=true,name=app_id"`
	// Search against name, app instance identifier, and app class ID.
	NameSearch *string `queryParam:"style=form,explode=true,name=name_search"`
	// Search filter should be an exact match.
	ExactMatch *bool `default:"false" queryParam:"style=form,explode=true,name=exact_match"`
	// Get all apps in the AppStore regardless of visibility. Only available to admins.
	AllVisibilities *bool  `default:"false" queryParam:"style=form,explode=true,name=all_visibilities"`
	Page            *int64 `default:"1" queryParam:"style=form,explode=true,name=page"`
	Size            *int64 `default:"50" queryParam:"style=form,explode=true,name=size"`
}

func (g GetAppStoreAppsRequest) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(g, "", false)
}

func (g *GetAppStoreAppsRequest) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &g, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *GetAppStoreAppsRequest) GetAppID() *string {
	if o == nil {
		return nil
	}
	return o.AppID
}

func (o *GetAppStoreAppsRequest) GetNameSearch() *string {
	if o == nil {
		return nil
	}
	return o.NameSearch
}

func (o *GetAppStoreAppsRequest) GetExactMatch() *bool {
	if o == nil {
		return nil
	}
	return o.ExactMatch
}

func (o *GetAppStoreAppsRequest) GetAllVisibilities() *bool {
	if o == nil {
		return nil
	}
	return o.AllVisibilities
}

func (o *GetAppStoreAppsRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *GetAppStoreAppsRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

type GetAppStoreAppsResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	PageAppStoreApp *shared.PageAppStoreApp
	// Validation Error
	HTTPValidationError *shared.HTTPValidationError
}

func (o *GetAppStoreAppsResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetAppStoreAppsResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetAppStoreAppsResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetAppStoreAppsResponse) GetPageAppStoreApp() *shared.PageAppStoreApp {
	if o == nil {
		return nil
	}
	return o.PageAppStoreApp
}

func (o *GetAppStoreAppsResponse) GetHTTPValidationError() *shared.HTTPValidationError {
	if o == nil {
		return nil
	}
	return o.HTTPValidationError
}
