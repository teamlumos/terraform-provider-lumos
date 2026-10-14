package hooks

import "net/http"

// appStoreSettingsExpand asks every app read for its full App Store
// settings, which lumos_app shows as a read-only nested block. expand is a
// list parameter, and the spec can't give it a default that the generated
// client sends.
type appStoreSettingsExpand struct{}

func (appStoreSettingsExpand) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if hookCtx.OperationID != "v1GetApp" {
		return req, nil
	}
	q := req.URL.Query()
	for _, v := range q["expand"] {
		if v == "appstore_settings" {
			return req, nil
		}
	}
	q.Add("expand", "appstore_settings")
	req.URL.RawQuery = q.Encode()
	return req, nil
}
