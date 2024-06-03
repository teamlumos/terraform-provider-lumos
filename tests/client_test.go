package provider

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/speakeasy/terraform-provider-lumos/internal/sdk"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/operations"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/models/shared"
)

const APP_ID = string("de97161b-13a8-7d16-8a06-722859c1b4a2")
const EMAIL = "albus@lumostester.com"
const USER_ID = "f7dae071-7640-5280-4764-0772993807ef"
const GROUP_ID = "b4aff9b6-f701-98a2-7404-d24a121ebf9c"
const GROUP_NAME = "OnlyAlbus"

func GetNewClient(t *testing.T) (*sdk.Lumos, context.Context) {
	apiToken := os.Getenv("LUMOS_API_TOKEN")
	if apiToken == "" {
		t.Fatalf(`LUMOS_API_TOKEN is not set`)
	}
	security := shared.Security{
		HTTPBearer: &apiToken,
	}
	opts := []sdk.SDKOption{
		sdk.WithServerURL("https://api.lumos.com"),
		sdk.WithSecurity(security),
		sdk.WithClient(http.DefaultClient),
	}
	client := sdk.New(opts...)
	return client, context.Background()
}

func TestGetApp(t *testing.T) {
	APP_NAME := string("Terraform Testing")
	client, ctx := GetNewClient(t)
	response, err := client.AppStore.GetAppStoreApp(ctx, operations.GetAppStoreAppRequest{
		AppID: APP_ID,
	})
	if err != nil || response == nil {
		t.Fatalf(`client.AppStore.GetAppStoreApp(%s) was unsuccessful, %v, %v`, APP_NAME, response, err)
	}
	if response.AppStoreApp.ID != APP_ID {
		t.Fatalf(`client.AppStore.GetAppStoreApp(%s) = %v, %v`, APP_NAME, response, err)
	}
}

// func TestSearchAppIsExactMatch(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.searchApp("Terraform Testin")
// 	if err == nil || response != nil {
// 		t.Fatalf(`searchApp("Terraform Testin") was successful but should not have been`)
// 	}
// }

// func TestGetApp(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.getApp(APP_ID)
// 	if err != nil || response == nil {
// 		t.Fatalf(`getApp(%s) was unsuccessful`, APP_ID)
// 	}
// 	if response.Id != APP_ID {
// 		t.Fatalf(`getApp(%s) = %v, %v`, APP_ID, response, err)
// 	}
// }

// func TestGetAppSettings(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.getAppSettings(APP_ID)
// 	if err != nil || response == nil {
// 		t.Fatalf(`getAppSetting(%s) was unsuccessful`, APP_ID)
// 	}
// }

// func TestAppCreateUpdate(t *testing.T) {
// 	client := GetNewClient(t)
// 	var app appResourceModel
// 	app.Name = types.StringValue("Terraform Test App - Create")
// 	app.Category = types.StringValue("Developers")
// 	app.Description = types.StringValue("Description")

// 	createResponse, err := client.createApp(app)
// 	if err != nil || createResponse == nil {
// 		t.Fatalf(`createApp(app) was unsuccessful, %v`, createResponse)
// 	}
// 	if createResponse.UserFriendlyLabel != "Terraform Test App - Create" {
// 		t.Fatalf(`createApp(app) = %v, %v`, createResponse, err)
// 	}

// 	app.Description = types.StringValue("Second description")
// 	app.Name = types.StringValue("Terraform Test App - Update")

// 	updateResponse, err := client.updateApp(createResponse.Id, app)
// 	if err != nil || updateResponse == nil {
// 		t.Fatalf(`updateApp(createResponse.Id, app) was unsuccessful, %v, %v`, updateResponse, err)
// 	}
// 	if updateResponse.UserFriendlyLabel != "Terraform Test App - Update" {
// 		t.Fatalf(`updateApp(createResponse.Id, app) = %v, %v`, updateResponse, err)
// 	}
// }

// func TestAddAppToAppStore(t *testing.T) {
// 	client := GetNewClient(t)

// 	app, err := client.createApp(appResourceModel{
// 		Name:        types.StringValue("Terraform Test App - Add"),
// 		Category:    types.StringValue("Developers"),
// 		Description: types.StringValue("Description"),
// 	})

// 	var group lumosAPIBaseObject
// 	group.Id = GROUP_ID

// 	var user lumosAPIBaseObject
// 	user.Id = USER_ID
// 	var appSettings appStoreSettingResourceModel
// 	settings, err := client.addAppToAppStore(app.Id, appSettings)
// 	if err == nil || settings != nil {
// 		t.Fatalf(`addAppToAppStore(%s) was successful but should not have been`, app.Id)
// 	}

// 	getApp, err := client.getApp(app.Id)
// 	if err == nil || getApp != nil {
// 		t.Fatalf(`getApp(%s) was successful but should not have been`, app.Id)
// 	}

// 	appSettings.RequestFlow.Admins.Groups = []lumosAPIBaseObject{group}
// 	appSettings.RequestFlow.Admins.Users = []lumosAPIBaseObject{user}
// 	appSettings.RequestFlow.Approvers.Groups = []lumosAPIBaseObject{group}
// 	appSettings.RequestFlow.Approvers.Users = []lumosAPIBaseObject{user}
// 	_, err = client.addAppToAppStore(app.Id, appSettings)
// 	if err != nil {
// 		t.Fatalf(`addAppToAppStore(%s, appSettings) was unsuccessful, %v`, app.Id, err)
// 	}
// 	app, err = client.getApp(app.Id)
// 	if err != nil || app == nil {
// 		t.Fatalf(`getApp(%s) was successful`, app.Id)
// 	}

// 	_, err = client.removeAppFromAppStore(app.Id, appSettings)
// 	if err != nil {
// 		t.Fatalf(`removeAppFromAppStore(%s, appSettings) was unsuccessful, %v`, app.Id, err)
// 	}

// 	getApp, err = client.getApp(app.Id)
// 	if err == nil || getApp != nil {
// 		t.Fatalf(`getApp(%s) was successful but should not have been`, app.Id)
// 	}
// }

// func TestSearchUser(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.searchUser(EMAIL)
// 	if err != nil || response == nil {
// 		t.Fatalf(`searchUser(%s) was unsuccessful`, EMAIL)
// 	}
// 	if response.Id != USER_ID {
// 		t.Fatalf(`searchUser(%s) = %v, %v`, EMAIL, response, err)
// 	}
// }

// func TestSearchUserIsExactMatch(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.searchUser("albus@lumostester.co")
// 	if err == nil || response != nil {
// 		t.Fatalf(`searchUser("albus@lumostester.co") was successful but should not have been`)
// 	}
// }

// func TestGetUser(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.getUser(USER_ID)
// 	if err != nil || response == nil {
// 		t.Fatalf(`getUser(%s) was unsuccessful`, USER_ID)
// 	}
// 	if response.Id != USER_ID {
// 		t.Fatalf(`getUser(%s) = %v, %v`, USER_ID, response, err)
// 	}
// }

// func TestGetGroup(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.getGroup(GROUP_ID)
// 	if err != nil || response == nil {
// 		t.Fatalf(`getGroup(%s) was unsuccessful`, GROUP_ID)
// 	}
// 	if response.Id != GROUP_ID {
// 		t.Fatalf(`getGroup(%s) = %v, %v`, GROUP_ID, response, err)
// 	}
// }

// func TestSearchGroup(t *testing.T) {
// 	client := GetNewClient(t)
// 	response, err := client.searchGroup(GROUP_NAME)
// 	if err != nil || response == nil {
// 		t.Fatalf(`searchGroup(%s) was unsuccessful`, GROUP_NAME)
// 	}
// 	if response.Id != GROUP_ID {
// 		t.Fatalf(`searchGroup(%s) = %v, %v`, GROUP_NAME, response, err)
// 	}

// }

// func TestPermissionCRUD(t *testing.T) {
// 	client := GetNewClient(t)
// 	var permission requestablePermissionResourceModel
// 	permission.AppId = types.StringValue(APP_ID)
// 	permission.Label = types.StringValue("Test Permission 2")

// 	createResponse, err := client.createPermission(permission)
// 	if err != nil || createResponse == nil {
// 		t.Fatalf(`createPermission(permission) was unsuccessful`)
// 	}
// 	if createResponse.Label != "Test Permission 2" {
// 		t.Fatalf(`createPermission(permission) = %v, %v`, createResponse, err)
// 	}

// 	permission.ManualInstructions = types.StringValue("Manual Instructions")
// 	updateResponse, err := client.updatePermission(createResponse.Id, permission)
// 	if err != nil || updateResponse == nil {
// 		t.Fatalf(`updatePermission(createResponse.Id, permission) was unsuccessful`)
// 	}
// 	if updateResponse.RequestConfig.RequestFulfillmentConfig.ManualInstructions != "Manual Instructions" {
// 		t.Fatalf(`updatePermission(createResponse.Id, permission) = %v, %v`, updateResponse, err)
// 	}

// 	err = client.deletePermission(createResponse.Id)
// 	if err != nil {
// 		t.Fatalf(`deletePermission(createResponse.Id) was unsuccessful`)
// 	}
// }
