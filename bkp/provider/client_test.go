package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const APP_ID = "de97161b-13a8-7d16-8a06-722859c1b4a2"
const APP_NAME = "Terraform Testing"
const EMAIL = "albus@lumostester.com"
const USER_ID = "f7dae071-7640-5280-4764-0772993807ef"
const GROUP_ID = "b4aff9b6-f701-98a2-7404-d24a121ebf9c"
const GROUP_NAME = "OnlyAlbus"

func GetNewClient(t *testing.T) *LumosAPIClient {
	apiToken := os.Getenv("LUMOS_API_TOKEN")
	if apiToken == "" {
		t.Fatalf(`LUMOS_API_TOKEN is not set`)
	}
	client, _ := NewLumosAPIClient("https://api.lumos.com", apiToken)
	return client
}

func TestSearchApp(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchApp(APP_NAME)
	if err != nil || response == nil {
		t.Fatalf(`searchApp(%s) was unsuccessful`, APP_NAME)
	}
	if response.Id != APP_ID {
		t.Fatalf(`searchApp(%s) = %v, %v`, APP_NAME, response, err)
	}
}

func TestSearchAppIsExactMatch(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchApp("Terraform Testin")
	if err == nil || response != nil {
		t.Fatalf(`searchApp("Terraform Testin") was successful but should not have been`)
	}
}

func TestGetApp(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.getApp(APP_ID)
	if err != nil || response == nil {
		t.Fatalf(`getApp(%s) was unsuccessful`, APP_ID)
	}
	if response.Id != APP_ID {
		t.Fatalf(`getApp(%s) = %v, %v`, APP_ID, response, err)
	}
}

func TestSearchUser(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchUser(EMAIL)
	if err != nil || response == nil {
		t.Fatalf(`searchUser(%s) was unsuccessful`, EMAIL)
	}
	if response.Id != USER_ID {
		t.Fatalf(`searchUser(%s) = %v, %v`, EMAIL, response, err)
	}
}

func TestSearchUserIsExactMatch(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchUser("albus@lumostester.co")
	if err == nil || response != nil {
		t.Fatalf(`searchUser("albus@lumostester.co") was successful but should not have been`)
	}
}

func TestGetUser(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.getUser(USER_ID)
	if err != nil || response == nil {
		t.Fatalf(`getUser(%s) was unsuccessful`, USER_ID)
	}
	if response.Id != USER_ID {
		t.Fatalf(`getUser(%s) = %v, %v`, USER_ID, response, err)
	}
}

func TestGetGroup(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.getGroup(GROUP_ID)
	if err != nil || response == nil {
		t.Fatalf(`getGroup(%s) was unsuccessful`, GROUP_ID)
	}
	if response.Id != GROUP_ID {
		t.Fatalf(`getGroup(%s) = %v, %v`, GROUP_ID, response, err)
	}
}

func TestSearchGroup(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchGroup(GROUP_NAME)
	if err != nil || response == nil {
		t.Fatalf(`searchGroup(%s) was unsuccessful`, GROUP_NAME)
	}
	if response.Id != GROUP_ID {
		t.Fatalf(`searchGroup(%s) = %v, %v`, GROUP_NAME, response, err)
	}

}

func TestPermissionCRUD(t *testing.T) {
	client := GetNewClient(t)
	var permission requestablePermissionResourceModel
	permission.AppId = types.StringValue(APP_ID)
	permission.Label = types.StringValue("Test Permission 2")

	createResponse, err := client.createPermission(permission)
	if err != nil || createResponse == nil {
		t.Fatalf(`createPermission(permission) was unsuccessful`)
	}
	if createResponse.Label != "Test Permission 2" {
		t.Fatalf(`createPermission(permission) = %v, %v`, createResponse, err)
	}

	permission.ManualInstructions = types.StringValue("Manual Instructions")
	updateResponse, err := client.updatePermission(createResponse.Id, permission)
	if err != nil || updateResponse == nil {
		t.Fatalf(`updatePermission(createResponse.Id, permission) was unsuccessful`)
	}
	if updateResponse.RequestConfig.RequestFulfillmentConfig.ManualInstructions != "Manual Instructions" {
		t.Fatalf(`updatePermission(createResponse.Id, permission) = %v, %v`, updateResponse, err)
	}

	err = client.deletePermission(createResponse.Id)
	if err != nil {
		t.Fatalf(`deletePermission(createResponse.Id) was unsuccessful`)
	}
}
