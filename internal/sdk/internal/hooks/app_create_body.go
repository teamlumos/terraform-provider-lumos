package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// appCreateBody drops app_class_id = "custom" from a lumos_app create. A
// custom app reads back with that value, so a config may set it, but the API
// takes any app_class_id on create as an integration to connect.
type appCreateBody struct{}

func (appCreateBody) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if hookCtx.OperationID != "v1CreateApp" || req.Body == nil {
		return req, nil
	}
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading app create body: %w", err)
	}
	body := raw
	fields := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &fields); err != nil {
			return nil, fmt.Errorf("decoding app create body: %w", err)
		}
	}
	var appClassID string
	if v, ok := fields["app_class_id"]; ok {
		_ = json.Unmarshal(v, &appClassID)
	}
	if appClassID == "custom" {
		delete(fields, "app_class_id")
		if body, err = json.Marshal(fields); err != nil {
			return nil, fmt.Errorf("encoding app create body: %w", err)
		}
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	req.ContentLength = int64(len(body))
	return req, nil
}
