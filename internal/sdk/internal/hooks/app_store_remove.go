package hooks

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// removeFromAppStore handles destroying lumos_app_appstore_settings and
// lumos_entitlement_appstore_settings, which call provider-only operations.
// This hook runs on those operations and turns each into the request the
// API actually expects: a PATCH of the settings that takes the app or
// entitlement out of the App Store. Terraform treats removing the settings
// as a delete, and the API models it as that PATCH, so the hook translates
// one into the other.
type removeFromAppStore struct{}

var removeFromAppStoreBodies = map[string][]byte{
	"v1RemoveAppFromAppStore":         []byte(`{"in_app_store":false}`),
	"v1RemoveEntitlementFromAppStore": []byte(`{"visible_in_appstore":false}`),
}

func (removeFromAppStore) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	body, ok := removeFromAppStoreBodies[hookCtx.OperationID]
	if !ok {
		return req, nil
	}
	req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
	req.URL.RawPath = strings.TrimSuffix(req.URL.RawPath, "/")

	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
