package hooks

import "net/http"

// readExpand asks app and entitlement reads, and the app list, for the
// expansions the provider shows: full App Store settings on lumos_app,
// lumos_entitlement and lumos_apps, which show them as a read-only nested
// block, and an app's custom attributes. expand is a list parameter, and the
// spec can't give it a default that the generated client sends.
type readExpand struct{}

var readExpandOperations = map[string][]string{
	"v1GetApp":       {"appstore_settings", "custom_attributes"},
	"getEntitlement": {"appstore_settings"},
	"v1ListApps":     {"appstore_settings", "custom_attributes"},
}

func (readExpand) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	expansions, ok := readExpandOperations[hookCtx.OperationID]
	if !ok {
		return req, nil
	}
	q := req.URL.Query()
	requested := map[string]bool{}
	for _, v := range q["expand"] {
		requested[v] = true
	}
	for _, e := range expansions {
		if !requested[e] {
			q.Add("expand", e)
		}
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}
