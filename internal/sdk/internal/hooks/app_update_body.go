package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// appUpdateBody shapes a lumos_app update so it follows the config. Every
// update sends each field the plan holds, so a missing field was left out
// of the config.
//
//   - app_class_id is dropped. It never changes on update, and the API
//     compares it with its own ID for the app, which isn't "custom" for a
//     custom app.
//   - null is sent for each *_override field and for post_approval_message
//     when it's missing, which clears it. The generated client drops nulls.
type appUpdateBody struct{}

var clearableAppMetadata = []string{
	"label_override",
	"description_override",
	"category_override",
	"logo_url_override",
	"website_url_override",
	"post_approval_message",
}

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

	delete(fields, "app_class_id")
	for _, name := range clearableAppMetadata {
		if _, ok := fields[name]; !ok {
			fields[name] = json.RawMessage("null")
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
