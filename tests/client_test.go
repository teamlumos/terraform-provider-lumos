package tests

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/teamlumos/terraform-provider-lumos/internal/sdk"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/operations"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
)

func GetNewClient(t *testing.T) (*sdk.Lumos, context.Context) {
	apiToken := os.Getenv("LUMOS_API_TOKEN")
	if apiToken == "" {
		t.Fatalf(`LUMOS_API_TOKEN is not set`)
	}
	security := shared.Security{
		HTTPBearer: apiToken,
	}
	opts := []sdk.SDKOption{
		sdk.WithServerURL(os.Getenv("LUMOS_API_URL")),
		sdk.WithSecurity(security),
		sdk.WithClient(http.DefaultClient),
	}
	client := sdk.New(opts...)
	return client, context.Background()
}

func TestListIntegrations(t *testing.T) {
	client, ctx := GetNewClient(t)
	response, err := client.Integrations.ListIntegrationsV1(ctx, operations.ListIntegrationsV1Request{})
	if err != nil || response == nil {
		t.Fatalf(`client.Integrations.ListIntegrationsV1() was unsuccessful, %v, %v`, response, err)
	}
	if response.IntegrationCatalogListOutput == nil {
		t.Fatalf(`client.Integrations.ListIntegrationsV1() returned no page payload, %v`, response)
	}
	if len(response.IntegrationCatalogListOutput.Integrations) == 0 {
		t.Fatalf(`client.Integrations.ListIntegrationsV1() returned zero integrations`)
	}
}

func TestGetApp(t *testing.T) {
	client, ctx := GetNewClient(t)
	appID := os.Getenv("APP_ID")
	if appID == "" {
		t.Skip(`APP_ID is not set`)
	}
	response, err := client.Apps.V1GetApp(ctx, operations.V1GetAppRequest{
		ID: appID,
	})
	if err != nil || response == nil {
		t.Fatalf(`client.Apps.V1GetApp(%s) was unsuccessful, %v, %v`, appID, response, err)
	}
	if response.AppResponse == nil {
		t.Fatalf(`client.Apps.V1GetApp(%s) returned no app, %v`, appID, response)
	}
	if response.AppResponse.ID != appID {
		t.Fatalf(`client.Apps.V1GetApp(%s) returned id %q`, appID, response.AppResponse.ID)
	}
}
