package hooks

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// removeAppFromAppStore handles destroying lumos_app_appstore_settings,
// which calls the provider-only v1RemoveAppFromAppStore operation. This hook
// runs on that operation and turns it into the request the API actually
// expects: a PATCH of the app's settings with in_app_store set to false.
// Terraform treats removing the settings as a delete, and the API models it
// as that PATCH, so the hook translates one into the other.
type removeAppFromAppStore struct{}

func (removeAppFromAppStore) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if hookCtx.OperationID != "v1RemoveAppFromAppStore" {
		return req, nil
	}
	req.URL.Path = strings.TrimSuffix(req.URL.Path, "/")
	req.URL.RawPath = strings.TrimSuffix(req.URL.RawPath, "/")

	body := []byte(`{"in_app_store":false}`)
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
