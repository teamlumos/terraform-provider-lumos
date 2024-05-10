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
var _ datasource.DataSource = &requestablePermissionDataSource{}

func NewRequestablePermissionDataSource() datasource.DataSource {
	return &requestablePermissionDataSource{}
}

// ExampleDataSource defines the data source implementation.
type requestablePermissionDataSource struct {
	client *LumosAPIClient
}

// requestablePermissionDataSourceModel describes the data source data model.
type requestablePermissionDataSourceModel struct {
	Id    types.String `tfsdk:"id"`
	AppId types.String `tfsdk:"app_id"`
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
			"app_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The app ID of this requestable permission.",
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

	client, ok := req.ProviderData.(*LumosAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *LumosAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *requestablePermissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state requestablePermissionDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	///////////////////// API Call
	permission, err := d.client.getPermission(state.Id.ValueString())
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	state.Id = types.StringValue(permission.Id)
	state.AppId = types.StringValue(permission.AppId)
	state.Label = types.StringValue(permission.Label)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
