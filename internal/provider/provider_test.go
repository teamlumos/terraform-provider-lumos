package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/operations"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func providerConfig() string {
	return fmt.Sprintf(`
provider "lumos" {
  http_bearer    = "%s"
  server_url     = "%s"
}
`, os.Getenv("LUMOS_API_TOKEN"), os.Getenv("LUMOS_API_URL"))
}

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"lumos": providerserver.NewProtocol6WithError(New("test")()),
	}
)

// testAccClient skips the test outside acceptance runs, matching resource.Test,
// since its callers run before resource.Test gets the chance to.
func testAccClient(t *testing.T) *sdk.Lumos {
	t.Helper()
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}
	return sdk.New(
		sdk.WithServerURL(os.Getenv("LUMOS_API_URL")),
		sdk.WithSecurity(shared.Security{HTTPBearer: os.Getenv("LUMOS_API_TOKEN")}),
		sdk.WithClient(http.DefaultClient),
	)
}

// expectNotUnknownValue is the inverse of plancheck.ExpectUnknownValue: it fails
// only when the planned value is "(known after apply)", so a value planned as
// known null still passes.
func expectNotUnknownValue(resourceAddress string, attributePath tfjsonpath.Path) plancheck.PlanCheck {
	return expectNotUnknownValueCheck{resourceAddress: resourceAddress, attributePath: attributePath}
}

type expectNotUnknownValueCheck struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

func (e expectNotUnknownValueCheck) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	for _, rc := range req.Plan.ResourceChanges {
		if rc.Address != e.resourceAddress {
			continue
		}
		result, err := tfjsonpath.Traverse(rc.Change.AfterUnknown, e.attributePath)
		if err != nil {
			// after_unknown only lists unknown values, so a missing path is known.
			return
		}
		if isUnknown, ok := result.(bool); ok && isUnknown {
			resp.Error = fmt.Errorf("Expected known value at %q, but found unknown value", e.attributePath.String())
		}
		return
	}
	resp.Error = fmt.Errorf("%s - Resource not found in plan ResourceChanges", e.resourceAddress)
}

// otherActiveUserID returns an active user in the test domain other than
// excludedID. It is resolved once, before any test step, so every step of a test
// sees the same user.
func otherActiveUserID(t *testing.T, excludedID string) string {
	t.Helper()
	res, err := testAccClient(t).Core.ListUsers(context.Background(), operations.ListUsersRequest{})
	if err != nil {
		t.Fatalf("listing users to pick a second approver: %v", err)
	}
	for _, user := range res.GetPageUserWithCustomAttributes().GetItems() {
		if user.ID != excludedID && user.Status != nil && *user.Status == shared.UserLifecycleStatusActive {
			t.Logf("using user %s as the second approver", user.ID)
			return user.ID
		}
	}
	t.Fatalf("no active user other than %s in the first page of users", excludedID)
	return ""
}

// otherGroupID returns a group in the test domain other than excludedID. It is
// resolved once, before any test step, so every step of a test sees the same
// group.
func otherGroupID(t *testing.T, excludedID string) string {
	t.Helper()
	res, err := testAccClient(t).Core.GetGroups(context.Background(), operations.GetGroupsRequest{})
	if err != nil {
		t.Fatalf("listing groups to pick a second approver group: %v", err)
	}
	for _, group := range res.GetPageGroup().GetItems() {
		if group.ID != nil && *group.ID != excludedID {
			t.Logf("using group %s as the second approver group", *group.ID)
			return *group.ID
		}
	}
	t.Fatalf("no group other than %s in the first page of groups", excludedID)
	return ""
}
