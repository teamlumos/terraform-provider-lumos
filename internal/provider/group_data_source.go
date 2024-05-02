// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &groupDataSource{}

func NewgroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

// groupDataSource defines the data source implementation.
type groupDataSource struct {
	apiToken string
}

// groupDataSourceModel describes the data source data model.
type groupDataSourceModel struct {
	Id                    types.String `tfsdk:"id"`
	AppId                 types.String `tfsdk:"app_id"`
	SourceAppId           types.String `tfsdk:"source_app_id"`
	IntegrationSpecificId types.String `tfsdk:"integration_specific_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Lifecycle             types.String `tfsdk:"lifecycle"`
}

func (d *groupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Group",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this group.",
			},
			"app_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the app that owns this group.",
			},
			"source_app_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the app that owns this group.",
			},
			"integration_specific_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of this group, specific to the integration.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of this group",
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of this group",
			},
			"lifecycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "An enumeration of SYNCED or NATIVE",
			},
		},
	}
}

func (d *groupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupDataSourceModel

	///////////////////// API Call

	// TODO: update to retrieve ID from Terraform configuration.
	client, _ := NewLumosAPIClient(BASE_URL, d.apiToken)
	state.Name = "Example Group"
	group, err := client.searchGroup(state.Name)
	if err != nil {
		fmt.Printf("Error getting group %s: %s", state.Name, err)
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	state.Id = group.Id
	state.Name = group.Name

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
