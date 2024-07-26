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
var _ datasource.DataSource = &RequestablePermissionsDataSource{}
var _ datasource.DataSourceWithConfigure = &RequestablePermissionsDataSource{}

func NewRequestablePermissionsDataSource() datasource.DataSource {
	return &RequestablePermissionsDataSource{}
}

// RequestablePermissionsDataSource is the data source implementation.
type RequestablePermissionsDataSource struct {
	client *sdk.Lumos
}

// RequestablePermissionsDataSourceModel describes the data model.
type RequestablePermissionsDataSourceModel struct {
	AppID      types.String                          `tfsdk:"app_id"`
	ExactMatch types.Bool                            `tfsdk:"exact_match"`
	InAppStore types.Bool                            `tfsdk:"in_app_store"`
	Items      []tfTypes.RequestablePermissionOutput `tfsdk:"items"`
	Page       types.Int64                           `tfsdk:"page"`
	Pages      types.Int64                           `tfsdk:"pages"`
	SearchTerm types.String                          `tfsdk:"search_term"`
	Size       types.Int64                           `tfsdk:"size"`
	Total      types.Int64                           `tfsdk:"total"`
}

// Metadata returns the data source type name.
func (r *RequestablePermissionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_requestable_permissions"
}

// Schema defines the schema for the data source.
func (r *RequestablePermissionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RequestablePermissions DataSource",

		Attributes: map[string]schema.Attribute{
			"app_id": schema.StringAttribute{
				Optional:    true,
				Description: `Filters requestable permissions by the ID of the app to which they belong.`,
			},
			"exact_match": schema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: `Search filter should be an exact match.`,
			},
			"in_app_store": schema.BoolAttribute{
				Optional: true,
				MarkdownDescription: `` + "\n" +
					`    Filters permissions by visibility in the AppStore.` + "\n" +
					``,
			},
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_class_id": schema.StringAttribute{
							Computed:    true,
							Description: `The non-unique ID of the service associated with this requestable permission. Depending on how it is sourced in Lumos, this may be the app's name, website,  or other identifier.`,
						},
						"app_id": schema.StringAttribute{
							Computed:    true,
							Description: `The ID of the app associated with this requestable permission.`,
						},
						"app_instance_id": schema.StringAttribute{
							Computed:    true,
							Description: `The ID of the instance associated with this requestable permission.`,
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Description: `The ID of this requestable permission.`,
						},
						"label": schema.StringAttribute{
							Computed:    true,
							Description: `The label of this requestable permission.`,
						},
						"request_config": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"access_removal_inline_webhook": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"description": schema.StringAttribute{
											Computed:    true,
											Description: `The description of this inline webhook.`,
										},
										"hook_type": schema.StringAttribute{
											Computed:    true,
											Description: `An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]`,
										},
										"id": schema.StringAttribute{
											Computed:    true,
											Description: `The ID of this inline webhook.`,
										},
										"name": schema.StringAttribute{
											Computed:    true,
											Description: `The name of this inline webhook.`,
										},
									},
									Description: `A deprovisioning webhook can be optionally associated with this config.`,
								},
								"allowed_groups": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"groups": schema.SetNestedAttribute{
											Computed: true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"app_id": schema.StringAttribute{
														Computed:    true,
														Description: `The ID of the app that sources this group.`,
													},
													"description": schema.StringAttribute{
														Computed:    true,
														Description: `The description of this group.`,
													},
													"group_lifecycle": schema.StringAttribute{
														Computed:    true,
														Description: `The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]`,
													},
													"id": schema.StringAttribute{
														Computed:    true,
														Description: `The ID of this group.`,
													},
													"integration_specific_id": schema.StringAttribute{
														Computed:    true,
														Description: `The ID of this group, specific to the integration.`,
													},
													"name": schema.StringAttribute{
														Computed:    true,
														Description: `The name of this group.`,
													},
													"source_app_id": schema.StringAttribute{
														Computed:    true,
														Description: `The ID of the app that sources this group.`,
													},
												},
											},
											Description: `The groups allowed to request this permission.`,
										},
										"type": schema.StringAttribute{
											Computed:    true,
											Description: `The type of this allowed groups config, can be all groups or specific. must be one of ["ALL_GROUPS", "SPECIFIED_GROUPS"]`,
										},
									},
									Description: `The allowed groups config associated with this config.`,
								},
								"allowed_groups_override": schema.BoolAttribute{
									Computed:    true,
									Description: `Indicates if allowed groups is overriden from the app-level settings.`,
								},
								"appstore_visibility": schema.StringAttribute{
									Computed:    true,
									Description: `The appstore visibility of this request config. must be one of ["HIDDEN", "VISIBLE"]`,
								},
								"request_approval_config": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"approvers": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"groups": schema.SetNestedAttribute{
													Computed: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"app_id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of the app that sources this group.`,
															},
															"description": schema.StringAttribute{
																Computed:    true,
																Description: `The description of this group.`,
															},
															"group_lifecycle": schema.StringAttribute{
																Computed:    true,
																Description: `The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]`,
															},
															"id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of this group.`,
															},
															"integration_specific_id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of this group, specific to the integration.`,
															},
															"name": schema.StringAttribute{
																Computed:    true,
																Description: `The name of this group.`,
															},
															"source_app_id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of the app that sources this group.`,
															},
														},
													},
													Description: `Groups assigned as support request approvers.`,
												},
												"users": schema.SetNestedAttribute{
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
													Description: `Users assigned as support request approvers.`,
												},
											},
											Description: `AppStore App approvers assigned.`,
										},
										"approvers_stage_2": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"groups": schema.SetNestedAttribute{
													Computed: true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"app_id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of the app that sources this group.`,
															},
															"description": schema.StringAttribute{
																Computed:    true,
																Description: `The description of this group.`,
															},
															"group_lifecycle": schema.StringAttribute{
																Computed:    true,
																Description: `The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]`,
															},
															"id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of this group.`,
															},
															"integration_specific_id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of this group, specific to the integration.`,
															},
															"name": schema.StringAttribute{
																Computed:    true,
																Description: `The name of this group.`,
															},
															"source_app_id": schema.StringAttribute{
																Computed:    true,
																Description: `The ID of the app that sources this group.`,
															},
														},
													},
													Description: `Groups assigned as support request approvers.`,
												},
												"users": schema.SetNestedAttribute{
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
													Description: `Users assigned as support request approvers.`,
												},
											},
											Description: `AppStore App stage 2 approvers assigned.`,
										},
										"custom_approval_message": schema.StringAttribute{
											Computed:    true,
											Description: `After the approval step, send a custom message to requesters. Note that the permission level approval message will override the App level approval message if custom_approval_message_override is set. Markdown for links and text formatting is supported.`,
										},
										"custom_approval_message_override": schema.BoolAttribute{
											Computed:    true,
											Description: `Indicates if custom_approval_message is overridden.`,
										},
										"manager_approval": schema.StringAttribute{
											Computed:    true,
											Description: `Manager approval can be configured as necessary to continue. must be one of ["NONE", "INITIAL_APPROVAL"]`,
										},
										"request_approval_config_override": schema.BoolAttribute{
											Computed:    true,
											Description: `Indicates if approval flow is overridden.`,
										},
										"require_additional_approval": schema.BoolAttribute{
											Computed:    true,
											Description: `Only turn on when working with sensitive permissions to ensure a smooth employee experience.`,
										},
									},
									Description: `A request approval config can be optionally associated with this config`,
								},
								"request_fulfillment_config": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"manual_instructions": schema.StringAttribute{
											Computed:    true,
											Description: `The manual instructions that go along.`,
										},
										"manual_steps_needed": schema.BoolAttribute{
											Computed:    true,
											Description: `Whether manual steps are needed.`,
										},
										"provisioning_group": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"app_id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of the app that sources this group.`,
												},
												"description": schema.StringAttribute{
													Computed:    true,
													Description: `The description of this group.`,
												},
												"group_lifecycle": schema.StringAttribute{
													Computed:    true,
													Description: `The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]`,
												},
												"id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of this group.`,
												},
												"integration_specific_id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of this group, specific to the integration.`,
												},
												"name": schema.StringAttribute{
													Computed:    true,
													Description: `The name of this group.`,
												},
												"source_app_id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of the app that sources this group.`,
												},
											},
											Description: `The provisioning group optionally assocated with this config.`,
										},
										"provisioning_webhook": schema.SingleNestedAttribute{
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"description": schema.StringAttribute{
													Computed:    true,
													Description: `The description of this inline webhook.`,
												},
												"hook_type": schema.StringAttribute{
													Computed:    true,
													Description: `An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]`,
												},
												"id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of this inline webhook.`,
												},
												"name": schema.StringAttribute{
													Computed:    true,
													Description: `The name of this inline webhook.`,
												},
											},
											Description: `The provisioning webhook optionally associated with this config.`,
										},
										"time_based_access": schema.ListAttribute{
											Computed:    true,
											ElementType: types.StringType,
											Description: `If enabled, users can request an app for a selected duration. After expiry, Lumos will automatically remove user's access.`,
										},
										"time_based_access_override": schema.BoolAttribute{
											Computed:    true,
											Description: `Indicates if time based access is overriden.`,
										},
									},
									Description: `A request fulfillment config can be optionally associated with this config`,
								},
								"request_validation_inline_webhook": schema.SingleNestedAttribute{
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"description": schema.StringAttribute{
											Computed:    true,
											Description: `The description of this inline webhook.`,
										},
										"hook_type": schema.StringAttribute{
											Computed:    true,
											Description: `An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]`,
										},
										"id": schema.StringAttribute{
											Computed:    true,
											Description: `The ID of this inline webhook.`,
										},
										"name": schema.StringAttribute{
											Computed:    true,
											Description: `The name of this inline webhook.`,
										},
									},
									Description: `A request validation webhook can be optionally associated with this config.`,
								},
							},
							Description: `The request config associated with this requestable permission.`,
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: `The type of this requestable permission. must be one of ["SYNCED", "NATIVE"]`,
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
				Optional: true,
				MarkdownDescription: `` + "\n" +
					`    Searches permissions by the permission's group name,` + "\n" +
					`    request configuration name, or specific integration ID.` + "\n" +
					``,
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

func (r *RequestablePermissionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *RequestablePermissionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *RequestablePermissionsDataSourceModel
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

	appID := new(string)
	if !data.AppID.IsUnknown() && !data.AppID.IsNull() {
		*appID = data.AppID.ValueString()
	} else {
		appID = nil
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
	inAppStore := new(bool)
	if !data.InAppStore.IsUnknown() && !data.InAppStore.IsNull() {
		*inAppStore = data.InAppStore.ValueBool()
	} else {
		inAppStore = nil
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
	request := operations.GetAppstorePermissionsForAppAppstoreRequestablePermissionsGetRequest{
		AppID:      appID,
		SearchTerm: searchTerm,
		ExactMatch: exactMatch,
		InAppStore: inAppStore,
		Page:       page,
		Size:       size,
	}
	res, err := r.client.AppStore.GetAppstorePermissionsForAppAppstoreRequestablePermissionsGet(ctx, request)
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
	if !(res.PageRequestablePermissionOutput != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedPageRequestablePermissionOutput(res.PageRequestablePermissionOutput)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
