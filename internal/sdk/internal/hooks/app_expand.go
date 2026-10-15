package hooks

import "net/http"

// appStoreSettingsExpand asks every app and entitlement read for its full
// App Store settings, which lumos_app and lumos_entitlement show as a
// read-only nested block. expand is a list parameter, and the spec can't
// give it a default that the generated client sends.
type appStoreSettingsExpand struct{}

var appStoreSettingsExpandOperations = map[string]bool{
	"v1GetApp":       true,
	"getEntitlement": true,
}

func (appStoreSettingsExpand) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if !appStoreSettingsExpandOperations[hookCtx.OperationID] {
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
