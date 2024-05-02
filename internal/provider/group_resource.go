// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &groupResource{}
var _ resource.ResourceWithImportState = &groupResource{}
var _ resource.ResourceWithImportState = &groupResource{}

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	client *LumosAPIClient
}

type groupResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	AppId                 types.String `tfsdk:"app_id"`
	SourceAppId           types.String `tfsdk:"source_app_id"`
	IntegrationSpecificId types.String `tfsdk:"integration_specific_id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	Lifecycle             types.String `tfsdk:"lifecycle"`
}

func (r *groupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Group",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of this group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the app that owns this group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"source_app_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of the app that owns this group.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_specific_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The ID of this group, specific to the integration.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of this group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The description of this group",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lifecycle": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "An enumeration of SYNCED or NATIVE",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *groupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	lumos_client, ok := req.ProviderData.(*LumosAPIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected string, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = lumos_client
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *groupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.getGroup(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting group",
			"Could not get group "+data.Id.ValueString()+" unexpected error:"+err.Error(),
		)
		return
	}

	*data = setUserResourceModelFromLumosAPIGroup(ctx, *data, *group)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
