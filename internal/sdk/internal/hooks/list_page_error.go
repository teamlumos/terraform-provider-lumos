package hooks

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// listPageError makes a list read fail when any of its pages comes back with
// an error status, so a data source never returns part of a list. The
// generated paging loop checks the status of the first page only, and the
// generated client returns error responses without an error, so this is the
// check for every later page.
type listPageError struct{}

var listPageErrorOperations = map[string]bool{
	"v1ListUsers":                    true,
	"v1ListApps":                     true,
	"listEntitlements":               true,
	"listAppEntitlements":            true,
	"listAccounts":                   true,
	"listGroups":                     true,
	"listInlineWebhooks":             true,
	"listAppCategories":              true,
	"listIntegrationsV1":             true,
	"listCustomAttributeDefinitions": true,
	"listWorkflowSteps":              true,
}

func (listPageError) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if !listPageErrorOperations[hookCtx.OperationID] || (res.StatusCode >= 200 && res.StatusCode < 300) {
		return res, nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
	return nil, fmt.Errorf("%s: status %d: %s", hookCtx.OperationID, res.StatusCode, bytes.TrimSpace(body))
}
