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
var _ datasource.DataSource = &appstoreAppDataSource{}

func NewAppstoreAppDataSource() datasource.DataSource {
	return &appstoreAppDataSource{}
}

// appstoreAppDataSource defines the data source implementation.
type appstoreAppDataSource struct {
	client *LumosAPIClient
}

// appstoreAppDataSourceModel describes the data source data model.
type appstoreAppDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	UserFriendlyLabel types.String `tfsdk:"user_friendly_label"`
	InstanceId        types.String `tfsdk:"instance_id"`
	AppClassId        types.String `tfsdk:"app_class_id"`
	Status            types.String `tfsdk:"status"`
}

func (d *appstoreAppDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_appstore_app"
}

func (d *appstoreAppDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "App",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this app.",
			},
			"app_class_id": schema.StringAttribute{
				Computed:    true,
				Description: "The app class ID of this app.",
			},
			"user_friendly_label": schema.StringAttribute{
				Required:    true,
				Description: "The user friendly label of this app.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The status of this app.",
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: "The instance ID of this app.",
			},
		},
	}
}

func (d *appstoreAppDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*LumosAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type in AppstoreAppDataSource",
			fmt.Sprintf("Expected LumosAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *appstoreAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state appstoreAppDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	appstoreApp, err := d.client.searchApp(state.UserFriendlyLabel.ValueString())
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	state.Id = types.StringValue(appstoreApp.Id)
	state.UserFriendlyLabel = types.StringValue(appstoreApp.UserFriendlyLabel)
	state.InstanceId = types.StringValue(appstoreApp.InstanceId)
	state.AppClassId = types.StringValue(appstoreApp.AppClassId)
	state.Status = types.StringValue(appstoreApp.Status)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a appstoreApp data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
