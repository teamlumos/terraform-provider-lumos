package hooks

import (
	"fmt"
	"net/http"
	"strings"
)

// requireListFilter fails a list of a large collection before it's sent
// when none of its narrowing filters is set, so a lookup can't page through
// every user, account, group or entitlement in the organization. The spec
// marks every filter optional, so the generated schema can't require one.
type requireListFilter struct{}

var requireListFilterOperations = map[string]struct {
	dataSource string
	filters    []string
}{
	"v1ListUsers":      {"lumos_users", []string{"search"}},
	"listEntitlements": {"lumos_entitlements", []string{"search"}},
	"listAccounts":     {"lumos_accounts", []string{"app_id", "email"}},
	"listGroups":       {"lumos_groups", []string{"name", "integration_specific_id", "app_id"}},
}

func (requireListFilter) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	op, ok := requireListFilterOperations[hookCtx.OperationID]
	if !ok {
		return req, nil
	}
	q := req.URL.Query()
	for _, name := range op.filters {
		if q.Get(name) != "" {
			return req, nil
		}
	}
	return nil, fmt.Errorf("%s requires %s", op.dataSource, joinOr(op.filters))
}

// joinOr renders ["a", "b", "c"] as "a, b or c".
func joinOr(names []string) string {
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " or " + names[len(names)-1]
}
