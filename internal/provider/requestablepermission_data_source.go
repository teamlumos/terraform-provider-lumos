// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

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
var _ datasource.DataSource = &RequestablePermissionDataSource{}
var _ datasource.DataSourceWithConfigure = &RequestablePermissionDataSource{}

func NewRequestablePermissionDataSource() datasource.DataSource {
	return &RequestablePermissionDataSource{}
}

// RequestablePermissionDataSource is the data source implementation.
type RequestablePermissionDataSource struct {
	client *sdk.Lumos
}

// RequestablePermissionDataSourceModel describes the data model.
type RequestablePermissionDataSourceModel struct {
	AppClassID    types.String                                    `tfsdk:"app_class_id"`
	AppID         types.String                                    `tfsdk:"app_id"`
	AppInstanceID types.String                                    `tfsdk:"app_instance_id"`
	ID            types.String                                    `tfsdk:"id"`
	Label         types.String                                    `tfsdk:"label"`
	RequestConfig tfTypes.RequestablePermissionInputRequestConfig `tfsdk:"request_config"`
	Type          types.String                                    `tfsdk:"type"`
}

// Metadata returns the data source type name.
func (r *RequestablePermissionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_requestable_permission"
}

// Schema defines the schema for the data source.
func (r *RequestablePermissionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "RequestablePermission DataSource",

		Attributes: map[string]schema.Attribute{
			"app_class_id": schema.StringAttribute{
				Computed:    true,
				Description: `The ID of the service associated with this requestable permission.`,
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
				Required: true,
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
						Description: `An inactivity workflow can be optionally associated with this config.`,
					},
					"allowed_groups": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"groups": schema.ListNestedAttribute{
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"app_id": schema.StringAttribute{
											Computed:    true,
											Description: `The ID of the app that owns this group.`,
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
											Description: `The ID of the app that owns this group.`,
										},
									},
								},
								Description: `The groups associated with this config.`,
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
						Description: `Indicates if allowed groups is overrided`,
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
									"groups": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"app_id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of the app that owns this group.`,
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
													Description: `The ID of the app that owns this group.`,
												},
											},
										},
										Description: `Groups assigned as support request approvers.`,
									},
									"users": schema.ListNestedAttribute{
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
									"groups": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"app_id": schema.StringAttribute{
													Computed:    true,
													Description: `The ID of the app that owns this group.`,
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
													Description: `The ID of the app that owns this group.`,
												},
											},
										},
										Description: `Groups assigned as support request approvers.`,
									},
									"users": schema.ListNestedAttribute{
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
								Description: `During the approval step, send a custom message to requesters. Note that the Permission level approval message will override the App level approval message. Markdown for links and text formatting is supported.`,
							},
							"custom_approval_message_override": schema.BoolAttribute{
								Computed:    true,
								Description: `Indicates if custom_approval_message is overrided`,
							},
							"manager_approval": schema.StringAttribute{
								Computed:    true,
								Description: `Manager approval can be configured as necessary to continue. must be one of ["NONE", "INITIAL_APPROVAL"]`,
							},
							"request_approval_config_override": schema.BoolAttribute{
								Computed:    true,
								Description: `Indicates if approval flow is overrided`,
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
										Description: `The ID of the app that owns this group.`,
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
										Description: `The ID of the app that owns this group.`,
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
								Description: `Indicates if time based access is overrided`,
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
				Description: `An enumeration. must be one of ["SYNCED", "NATIVE"]`,
			},
		},
	}
}

func (r *RequestablePermissionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *RequestablePermissionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *RequestablePermissionDataSourceModel
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

	permissionID := data.ID.ValueString()
	request := operations.GetAppstorePermissionAppstoreRequestablePermissionsPermissionIDGetRequest{
		PermissionID: permissionID,
	}
	res, err := r.client.AppStore.GetAppstorePermissionAppstoreRequestablePermissionsPermissionIDGet(ctx, request)
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
	if !(res.RequestablePermissionOutput != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedRequestablePermissionOutput(res.RequestablePermissionOutput)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
