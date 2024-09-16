// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/operations"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AppsDataSource{}
var _ datasource.DataSourceWithConfigure = &AppsDataSource{}

func NewAppsDataSource() datasource.DataSource {
	return &AppsDataSource{}
}

// AppsDataSource is the data source implementation.
type AppsDataSource struct {
	client *sdk.Lumos
}

// AppsDataSourceModel describes the data model.
type AppsDataSourceModel struct {
	ExactMatch types.Bool    `tfsdk:"exact_match"`
	Items      []tfTypes.App `tfsdk:"items"`
	NameSearch types.String  `tfsdk:"name_search"`
	Page       types.Int64   `tfsdk:"page"`
	Pages      types.Int64   `tfsdk:"pages"`
	Size       types.Int64   `tfsdk:"size"`
	Total      types.Int64   `tfsdk:"total"`
}

// Metadata returns the data source type name.
func (r *AppsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apps"
}

// Schema defines the schema for the data source.
func (r *AppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Apps DataSource",

		Attributes: map[string]schema.Attribute{
			"exact_match": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: `Search filter should be an exact match.`,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:    true,
							Description: `The ID of this app.`,
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
							Description: `The status of this app. Possible values: 'DISCOVERED', 'NEEDS_REVIEW', 'APPROVED', 'BLOCKLISTED', 'DEPRECATED'`,
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
				},
			},
			"name_search": schema.StringAttribute{
				Optional:    true,
				Description: `Search against name, app instance identifier, and app class ID.`,
			},
			"page": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: `Page number`,
			},
			"pages": schema.Int64Attribute{
				Computed: true,
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Optional:    true,
				Description: `Page size`,
			},
			"total": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *AppsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *AppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AppsDataSourceModel
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

	nameSearch := new(string)
	if !data.NameSearch.IsUnknown() && !data.NameSearch.IsNull() {
		*nameSearch = data.NameSearch.ValueString()
	} else {
		nameSearch = nil
	}
	exactMatch := new(bool)
	if !data.ExactMatch.IsUnknown() && !data.ExactMatch.IsNull() {
		*exactMatch = data.ExactMatch.ValueBool()
	} else {
		exactMatch = nil
	}
	page := new(int64)
	if !data.Page.IsUnknown() && !data.Page.IsNull() {
		*page = data.Page.ValueInt64()
	} else {
		page = nil
	}
	size := new(int64)
	if !data.Size.IsUnknown() && !data.Size.IsNull() {
		*size = data.Size.ValueInt64()
	} else {
		size = nil
	}
	request := operations.ListAppsRequest{
		NameSearch: nameSearch,
		ExactMatch: exactMatch,
		Page:       page,
		Size:       size,
	}
	res, err := r.client.Core.ListApps(ctx, request)
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
	if !(res.PageApp != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedPageApp(res.PageApp)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
