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

func TestGetApp(t *testing.T) {
	client, ctx := GetNewClient(t)
	appId := os.Getenv("APP_ID")
	response, err := client.AppStore.GetAppStoreApp(ctx, operations.GetAppStoreAppRequest{
		AppID: appId,
	})
	if err != nil || response == nil {
		t.Fatalf(`client.AppStore.GetAppStoreApp(%s) was unsuccessful, %v, %v`, appId, response, err)
	}
	if response.AppStoreApp.ID != appId {
		t.Fatalf(`client.AppStore.GetAppStoreApp(%s) = %v, %v`, appId, response, err)
	}
}

func TestCreateApp(t *testing.T) {
	client, ctx := GetNewClient(t)
	response, err := client.Core.CreateApp(ctx, shared.AppInputCreate{
		Name:        "Terraform Testing - Create app",
		Description: "Terraform Testing Description",
		Category:    "Developers",
	})
	if err != nil || response == nil {
		t.Fatalf(`client.Core.CreateApp() was unsuccessful, %v, %v`, response, err)
	}
}
