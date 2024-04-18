// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &requestablePermissionDataSource{}

func NewrequestablePermissionDataSource() datasource.DataSource {
	return &requestablePermissionDataSource{}
}

// ExampleDataSource defines the data source implementation.
type requestablePermissionDataSource struct {
	apiToken string
}

// requestablePermissionDataSourceModel describes the data source data model.
type requestablePermissionDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	Label types.String `tfsdk:"label"`
}

func (d *requestablePermissionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_requestable_permission"
}

func (d *requestablePermissionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Requestable Permission",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of this requestable permission.",
			},
			"label": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The label of this requestable permission.",
			},
		},
	}
}

func (d *requestablePermissionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiToken, ok := req.ProviderData.(string)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.apiToken = apiToken
}

func (d *requestablePermissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state requestablePermissionDataSourceModel

	///////////////////// API Call

	// TODO: update to retrieve ID from Terraform configuration.
	url := fmt.Sprintf("https://api.lumos.com/appstore/requestable_permissions/95997")
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	httpReq.Header.Set("Authorization", "Bearer "+d.apiToken)

	client := &http.Client{}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", httpResp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Assuming the API response is in JSON format
	var response struct {
		ID    string `json:"id"`
		Label string `json:"label"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("ID:", response.ID)
	fmt.Println("Label:", response.Label)
	///////////////////// API Call

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	state.Id = types.StringValue(response.ID)
	state.Label = types.StringValue(response.Label)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
