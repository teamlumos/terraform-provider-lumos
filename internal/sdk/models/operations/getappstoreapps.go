// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetAppStoreAppsRequest struct {
	// Filters apps by the ID of the service (i.e. okta.com). This parameter also requires app_instance_id to be included.
	AppClassID *string `queryParam:"style=form,explode=true,name=app_class_id"`
	// Filters apps by the ID of the instance from the service (app_class_id) from the service.This parameter also requires app_class_id to be included.
	AppInstanceID *string `queryParam:"style=form,explode=true,name=app_instance_id"`
	// Filters apps by the ID of the app.
	AppID *string `queryParam:"style=form,explode=true,name=app_id"`
	// Filters apps by name. This parameter requires a domain_id to be included.
	NameSearch *string `queryParam:"style=form,explode=true,name=name_search"`
	// Filters apps by name. This parameter requires a domain_id to be included.
	ExactMatch *bool  `default:"false" queryParam:"style=form,explode=true,name=exact_match"`
	Page       *int64 `default:"1" queryParam:"style=form,explode=true,name=page"`
	Size       *int64 `default:"50" queryParam:"style=form,explode=true,name=size"`
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

func (o *GetAppStoreAppsRequest) GetAppClassID() *string {
	if o == nil {
		return nil
	}
	return o.AppClassID
}

func (o *GetAppStoreAppsRequest) GetAppInstanceID() *string {
	if o == nil {
		return nil
	}
	return o.AppInstanceID
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
