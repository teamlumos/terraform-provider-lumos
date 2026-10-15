package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// inheritUnsetEntitlementOverrides makes an override that the
// lumos_entitlement_appstore_settings config leaves out mean "inherit the
// app's setting". The generated client omits unset fields, and the API's
// merge patch treats an omitted field as unchanged, so without this hook an
// override removed from the config would stay in place. The hook sends null
// for each override missing from the body, which the API reads as inherit.
//
// A missing key always means unset in config: the overrides aren't
// computed, and the settings PATCH sends every configured field, not only
// the changed ones. visible_in_appstore is left out on purpose; omitting it
// keeps its current value.
type inheritUnsetEntitlementOverrides struct{}

var entitlementOverrideFields = []string{
	"allowed_requester_group_uuids_override",
	"time_based_access_override",
	"request_validation_webhook_uuid",
}

func (inheritUnsetEntitlementOverrides) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	if hookCtx.OperationID != "updateEntitlementAppStoreSettings" {
		return req, nil
	}

	fields := map[string]json.RawMessage{}
	if req.Body != nil {
		raw, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("reading entitlement App Store settings body: %w", err)
		}
		if len(bytes.TrimSpace(raw)) > 0 {
			if err := json.Unmarshal(raw, &fields); err != nil {
				return nil, fmt.Errorf("decoding entitlement App Store settings body: %w", err)
			}
		}
	}
	for _, field := range entitlementOverrideFields {
		if _, ok := fields[field]; !ok {
			fields[field] = json.RawMessage("null")
		}
	}

	body, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("encoding entitlement App Store settings body: %w", err)
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
