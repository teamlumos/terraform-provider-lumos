package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	response, err := client.searchApp("Terraform Testing")
	if err != nil || response == nil {
		t.Fatalf(`searchApp("Terraform Testing") was unsuccessful`)
	}
	if response.Id != "de97161b-13a8-7d16-8a06-722859c1b4a2" {
		t.Fatalf(`searchApp("Terraform Testing") = %v, %v`, response, err)
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
	response, err := client.getApp("de97161b-13a8-7d16-8a06-722859c1b4a2")
	if err != nil || response == nil {
		t.Fatalf(`getApp("de97161b-13a8-7d16-8a06-722859c1b4a2") was unsuccessful`)
	}
	if response.Id != "de97161b-13a8-7d16-8a06-722859c1b4a2" {
		t.Fatalf(`getApp("de97161b-13a8-7d16-8a06-722859c1b4a2") = %v, %v`, response, err)
	}
}

func TestSearchUser(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchUser("albus@lumostester.com")
	if err != nil || response == nil {
		t.Fatalf(`searchUser("albus@lumostester.com") was unsuccessful`)
	}
	if response.Id != "f7dae071-7640-5280-4764-0772993807ef" {
		t.Fatalf(`searchUser("albus@lumostester.com") = %v, %v`, response, err)
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
	response, err := client.getUser("f7dae071-7640-5280-4764-0772993807ef")
	if err != nil || response == nil {
		t.Fatalf(`getUser("f7dae071-7640-5280-4764-0772993807ef") was unsuccessful`)
	}
	if response.Id != "f7dae071-7640-5280-4764-0772993807ef" {
		t.Fatalf(`getUser("f7dae071-7640-5280-4764-0772993807ef") = %v, %v`, response, err)
	}
}

func TestGetGroup(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.getGroup("b4aff9b6-f701-98a2-7404-d24a121ebf9c")
	if err != nil || response == nil {
		t.Fatalf(`getGroup("b4aff9b6-f701-98a2-7404-d24a121ebf9c") was unsuccessful`)
	}
	if response.Id != "b4aff9b6-f701-98a2-7404-d24a121ebf9c" {
		t.Fatalf(`getGroup("b4aff9b6-f701-98a2-7404-d24a121ebf9c") = %v, %v`, response, err)
	}
}

func TestSearchGroup(t *testing.T) {
	client := GetNewClient(t)
	response, err := client.searchGroup("OnlyAlbus")
	if err != nil || response == nil {
		t.Fatalf(`searchGroup("OnlyAlbus") was unsuccessful`)
	}
	if response.Id != "b4aff9b6-f701-98a2-7404-d24a121ebf9c" {
		t.Fatalf(`searchGroup("OnlyAlbus") = %v, %v`, response, err)
	}

}

func TestPermissionCRUD(t *testing.T) {
	client := GetNewClient(t)
	var permission requestablePermissionResourceModel
	permission.AppId = types.StringValue("de97161b-13a8-7d16-8a06-722859c1b4a2")
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
