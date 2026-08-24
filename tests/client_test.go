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
	response, err := client.Integrations.ListIntegrations(ctx, operations.ListIntegrationsRequest{})
	if err != nil || response == nil {
		t.Fatalf(`client.Integrations.ListIntegrations() was unsuccessful, %v, %v`, response, err)
	}
	if response.PageIntegrationCatalogOutput == nil {
		t.Fatalf(`client.Integrations.ListIntegrations() returned no page payload, %v`, response)
	}
	if len(response.PageIntegrationCatalogOutput.Items) == 0 {
		t.Fatalf(`client.Integrations.ListIntegrations() returned zero integrations`)
	}
}

func TestGetAppV1(t *testing.T) {
	client, ctx := GetNewClient(t)
	appID := os.Getenv("APP_ID")
	if appID == "" {
		t.Skip(`APP_ID is not set`)
	}
	response, err := client.Apps.GetAppV1(ctx, operations.GetAppV1Request{
		ID: appID,
	})
	if err != nil || response == nil {
		t.Fatalf(`client.Apps.GetAppV1(%s) was unsuccessful, %v, %v`, appID, response, err)
	}
	if response.AppWithCustomAttributes == nil {
		t.Fatalf(`client.Apps.GetAppV1(%s) returned no app, %v`, appID, response)
	}
	if response.AppWithCustomAttributes.ID != appID {
		t.Fatalf(`client.Apps.GetAppV1(%s) returned id %q`, appID, response.AppWithCustomAttributes.ID)
	}
}
