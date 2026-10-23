package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// appUpdateBody shapes a lumos_app update so it follows the config. Every
// update sends each field the plan holds, so a missing metadata field was
// left out of the config.
//
//   - app_class_id is dropped. It never changes on update, and the API
//     compares it with its own ID for the app, which isn't "custom" for a
//     custom app.
//   - For a custom app, null is sent for logo_url and website_url when
//     they're missing, which clears them. The generated client drops nulls.
//   - For a connected app, the metadata fields are dropped. The config
//     can't set them, so the plan holds the integration's values, and
//     sending those back would turn them into overrides.
type appUpdateBody struct{}

var (
	appMetadataFields          = []string{"label", "category", "description", "logo_url", "website_url"}
	customAppClearableMetadata = []string{"logo_url", "website_url"}
)

func (appUpdateBody) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if hookCtx.OperationID != "v1UpdateApp" || req.Body == nil {
		return req, nil
	}
	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("reading app update body: %w", err)
	}
	fields := map[string]json.RawMessage{}
	if len(bytes.TrimSpace(raw)) > 0 {
		if err := json.Unmarshal(raw, &fields); err != nil {
			return nil, fmt.Errorf("decoding app update body: %w", err)
		}
	}

	var appClassID string
	if v, ok := fields["app_class_id"]; ok {
		_ = json.Unmarshal(v, &appClassID)
	}
	delete(fields, "app_class_id")
	if appClassID == "" || appClassID == "custom" {
		for _, name := range customAppClearableMetadata {
			if _, ok := fields[name]; !ok {
				fields[name] = json.RawMessage("null")
			}
		}
	} else {
		for _, name := range appMetadataFields {
			delete(fields, name)
		}
	}

	body, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("encoding app update body: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
