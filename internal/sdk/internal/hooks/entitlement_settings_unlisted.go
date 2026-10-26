package hooks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// entitlementSettingsUnlisted makes the read of an entitlement's App Store
// settings fail, instead of answering 404, while the entitlement's app isn't
// in the App Store. The API answers 404 then although the settings still
// exist, and the provider takes a 404 as deleted outside Terraform: it drops
// lumos_entitlement_appstore_settings from state and plans a create that
// the API refuses until the app is listed again.
//
// The API also answers 404 when the entitlement is gone or has no settings,
// so on a 404 the hook reads the entitlement for its app, then that app's
// settings for in_app_store. The entitlement's own expanded
// appstore_settings can't tell the cases apart: the API leaves them out
// whenever it would answer this 404, whatever the reason. A 404 from either
// read, or an app in the App Store, leaves the 404 as it is. A read that
// fails any other way fails the settings read too, since the settings may
// still exist.
//
// It is an AfterSuccess hook because the client runs those for every
// response the server sends, whatever its status code.
type entitlementSettingsUnlisted struct{}

func (entitlementSettingsUnlisted) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if hookCtx.OperationID != "getEntitlementAppStoreSettings" || res.StatusCode != http.StatusNotFound || res.Request == nil {
		return res, nil
	}
	settingsURL := res.Request.URL
	base, rest, ok := strings.Cut(settingsURL.Path, "/v1/entitlements/")
	if !ok {
		return res, nil
	}
	entitlementID := strings.TrimSuffix(strings.TrimSuffix(rest, "/"), "/appstore-settings")

	var entitlement struct {
		AppID string `json:"app_id"`
	}
	found, err := entitlementSettingsLookup(hookCtx, res.Request, base+"/v1/entitlements/"+entitlementID, &entitlement)
	if err != nil {
		return nil, fmt.Errorf("entitlement %s: its App Store settings read back 404, and reading the entitlement to check its app failed: %w", entitlementID, err)
	}
	if !found {
		return res, nil
	}

	var app struct {
		InAppStore *bool `json:"in_app_store"`
	}
	found, err = entitlementSettingsLookup(hookCtx, res.Request, base+"/v1/apps/"+entitlement.AppID+"/appstore-settings", &app)
	if err != nil {
		return nil, fmt.Errorf("entitlement %s: its App Store settings read back 404, and reading its app's (%s) App Store settings failed: %w", entitlementID, entitlement.AppID, err)
	}
	if !found {
		return res, nil
	}
	if app.InAppStore == nil {
		return nil, fmt.Errorf("entitlement %s: its App Store settings read back 404, and its app's (%s) App Store settings have no in_app_store", entitlementID, entitlement.AppID)
	}
	if *app.InAppStore {
		return res, nil
	}
	// Removing the resource from the config would plan a destroy, which
	// refreshes it first and so fails here too; a removed block doesn't.
	res.Body.Close()
	return nil, fmt.Errorf("entitlement %s: the entitlement's app (%s) isn't in the App Store: list the app again, or stop managing this resource with a removed block that sets destroy = false", entitlementID, entitlement.AppID)
}

// entitlementSettingsLookup GETs path on the settings read's server with the
// same credentials, decoding a 200 into out. It reports false for a 404.
func entitlementSettingsLookup(hookCtx AfterSuccessContext, settingsReq *http.Request, path string, out any) (bool, error) {
	u := *settingsReq.URL
	u.Path, u.RawPath, u.RawQuery = path, "", ""
	req, err := http.NewRequestWithContext(hookCtx.Context, http.MethodGet, u.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", settingsReq.Header.Get("Authorization"))
	req.Header.Set("Accept", "application/json")
	got, err := hookCtx.SDKConfiguration.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer got.Body.Close()
	switch got.StatusCode {
	case http.StatusOK:
		return true, json.NewDecoder(got.Body).Decode(out)
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("status %d", got.StatusCode)
	}
}
