package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// clearUnsetFields makes a clearable App Store setting that the config
// leaves out mean "clear it". The generated client omits unset fields, and
// the API's merge patch treats an omitted field as unchanged, so without
// this hook a setting removed from the config would stay in place. The hook
// sends null for each clearable field missing from the body.
//
// A missing key always means unset in config: these fields aren't computed,
// and both settings PATCHes send every field the plan holds, not only the
// changed ones. Settings the API can't clear, such as visible_in_appstore
// and in_app_store, are required, and omitted lists plan as [], which the
// client sends, so neither is listed.
//
// A dotted name is a field of a nested object, cleared when the object is in
// the body. lumos_app_appstore_settings always sends time_based_access, since
// omitting it plans empty options and no default.
type clearUnsetFields struct{}

var clearUnsetFieldsByOperation = map[string][]string{
	// For an entitlement, a cleared override means it inherits the app's
	// setting.
	"updateEntitlementAppStoreSettings": {
		"allowed_requester_group_uuids_override",
		"time_based_access_override",
		"request_validation_webhook_uuid",
	},
	"v1PatchAppStoreSettings": {
		"request_instructions",
		"request_validation_webhook_uuid",
		"redirect_url",
		"default_entitlement",
		"time_based_access.default",
	},
}

func (clearUnsetFields) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	clearable, ok := clearUnsetFieldsByOperation[hookCtx.OperationID]
	if !ok {
		return req, nil
	}

	fields := map[string]json.RawMessage{}
	if req.Body != nil {
		raw, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading %s body: %w", hookCtx.OperationID, err)
		}
		if len(bytes.TrimSpace(raw)) > 0 {
			if err := json.Unmarshal(raw, &fields); err != nil {
				return nil, fmt.Errorf("decoding %s body: %w", hookCtx.OperationID, err)
			}
		}
	}
	for _, name := range clearable {
		parent, child, nested := strings.Cut(name, ".")
		if !nested {
			if _, ok := fields[name]; !ok {
				fields[name] = json.RawMessage("null")
			}
			continue
		}
		raw, ok := fields[parent]
		if !ok || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
			continue
		}
		object := map[string]json.RawMessage{}
		if err := json.Unmarshal(raw, &object); err != nil {
			return nil, fmt.Errorf("decoding %s %s: %w", hookCtx.OperationID, parent, err)
		}
		if _, ok := object[child]; ok {
			continue
		}
		object[child] = json.RawMessage("null")
		encoded, err := json.Marshal(object)
		if err != nil {
			return nil, fmt.Errorf("encoding %s %s: %w", hookCtx.OperationID, parent, err)
		}
		fields[parent] = encoded
	}

	body, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("encoding %s body: %w", hookCtx.OperationID, err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
