package hooks

import "net/http"

// appStoreSettingsExpand asks app reads for their full App Store settings,
// which the apps lookup shows on each app. expand is a list parameter, and
// the spec can't give it a default that the generated client sends.
type appStoreSettingsExpand struct{}

var appStoreSettingsExpandOperations = map[string]bool{
	"v1ListApps": true,
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
