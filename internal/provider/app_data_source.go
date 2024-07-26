// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/operations"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AppDataSource{}
var _ datasource.DataSourceWithConfigure = &AppDataSource{}

func NewAppDataSource() datasource.DataSource {
	return &AppDataSource{}
}

// AppDataSource is the data source implementation.
type AppDataSource struct {
	client *sdk.Lumos
}

// AppDataSourceModel describes the data model.
type AppDataSourceModel struct {
	AllowMultiplePermissionSelection types.Bool     `tfsdk:"allow_multiple_permission_selection"`
	AppClassID                       types.String   `tfsdk:"app_class_id"`
	ID                               types.String   `tfsdk:"id"`
	InstanceID                       types.String   `tfsdk:"instance_id"`
	LogoURL                          types.String   `tfsdk:"logo_url"`
	RequestInstructions              types.String   `tfsdk:"request_instructions"`
	Sources                          []types.String `tfsdk:"sources"`
	Status                           types.String   `tfsdk:"status"`
	UserFriendlyLabel                types.String   `tfsdk:"user_friendly_label"`
	WebsiteURL                       types.String   `tfsdk:"website_url"`
}

// Metadata returns the data source type name.
func (r *AppDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

// Schema defines the schema for the data source.
func (r *AppDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "App DataSource",

		Attributes: map[string]schema.Attribute{
			"allow_multiple_permission_selection": schema.BoolAttribute{
				Computed:    true,
				Description: `Determines whether users can request multiple permissions at once.This field will be removed in subsequent API versions.`,
			},
			"app_class_id": schema.StringAttribute{
				Computed:    true,
				Description: `The non-unique ID of the service associated with this requestable permission. Depending on how it is sourced in Lumos, this may be the app's name, website, or other identifier.`,
			},
			"id": schema.StringAttribute{
				Required: true,
			},
			"instance_id": schema.StringAttribute{
				Computed:    true,
				Description: `The non-unique ID of the instance associated with this app. This will be the Okta app id if it’s an Okta app, or will be marked as custom_app_import if manually uploaded into Lumos.`,
			},
			"logo_url": schema.StringAttribute{
				Computed:    true,
				Description: `The URL of the logo of this app.`,
			},
			"request_instructions": schema.StringAttribute{
				Computed:    true,
				Description: `The request instructions.`,
			},
			"sources": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: `The sources of this app.`,
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: `An enumeration. must be one of ["DISCOVERED", "NEEDS_REVIEW", "APPROVED", "BLOCKLISTED", "DEPRECATED"]`,
			},
			"user_friendly_label": schema.StringAttribute{
				Computed:    true,
				Description: `The user-friendly label of this app.`,
			},
			"website_url": schema.StringAttribute{
				Computed:    true,
				Description: `The URL of the website of this app.`,
			},
		},
	}
}

func (r *AppDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.Lumos)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *sdk.Lumos, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *AppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AppDataSourceModel
	var item types.Object

	resp.Diagnostics.Append(req.Config.Get(ctx, &item)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(item.As(ctx, &data, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()
	request := operations.GetAppRequest{
		ID: id,
	}
	res, err := r.client.Core.GetApp(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("failure to invoke API", err.Error())
		if res != nil && res.RawResponse != nil {
			resp.Diagnostics.AddError("unexpected http request/response", debugResponse(res.RawResponse))
		}
		return
	}
	if res == nil {
		resp.Diagnostics.AddError("unexpected response from API", fmt.Sprintf("%v", res))
		return
	}
	if res.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}
	if res.StatusCode != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", res.StatusCode), debugResponse(res.RawResponse))
		return
	}
	if !(res.App != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedApp(res.App)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
