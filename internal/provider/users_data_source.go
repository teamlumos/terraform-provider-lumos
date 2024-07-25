// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

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
var _ datasource.DataSource = &UsersDataSource{}
var _ datasource.DataSourceWithConfigure = &UsersDataSource{}

func NewUsersDataSource() datasource.DataSource {
	return &UsersDataSource{}
}

// UsersDataSource is the data source implementation.
type UsersDataSource struct {
	client *sdk.Lumos
}

// UsersDataSourceModel describes the data model.
type UsersDataSourceModel struct {
	ExactMatch types.Bool     `tfsdk:"exact_match"`
	Items      []tfTypes.User `tfsdk:"items"`
	Page       types.Int64    `tfsdk:"page"`
	Pages      types.Int64    `tfsdk:"pages"`
	SearchTerm types.String   `tfsdk:"search_term"`
	Size       types.Int64    `tfsdk:"size"`
	Total      types.Int64    `tfsdk:"total"`
}

// Metadata returns the data source type name.
func (r *UsersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (r *UsersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Users DataSource",

		Attributes: map[string]schema.Attribute{
			"exact_match": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: `If a search_term is provided, only accept exact matches.`,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"email": schema.StringAttribute{
							Computed:    true,
							Description: `The email of this user.`,
						},
						"family_name": schema.StringAttribute{
							Computed:    true,
							Description: `The family name of this user.`,
						},
						"given_name": schema.StringAttribute{
							Computed:    true,
							Description: `The given name of this user.`,
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: `The ID of this user.`,
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: `An enumeration. must be one of ["STAGED", "ACTIVE", "SUSPENDED", "INACTIVE"]`,
						},
					},
				},
			},
			"page": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"pages": schema.Int64Attribute{
				Computed: true,
			},
			"search_term": schema.StringAttribute{
				Optional:    true,
				Description: `Search for users by name or email.`,
			},
			"size": schema.Int64Attribute{
				Computed: true,
				Optional: true,
			},
			"total": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (r *UsersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *UsersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *UsersDataSourceModel
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

	searchTerm := new(string)
	if !data.SearchTerm.IsUnknown() && !data.SearchTerm.IsNull() {
		*searchTerm = data.SearchTerm.ValueString()
	} else {
		searchTerm = nil
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
	request := operations.ListUsersRequest{
		SearchTerm: searchTerm,
		ExactMatch: exactMatch,
		Page:       page,
		Size:       size,
	}
	res, err := r.client.Core.ListUsers(ctx, request)
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
	if !(res.PageUser != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedPageUser(res.PageUser)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
