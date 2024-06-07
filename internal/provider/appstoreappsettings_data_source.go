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
var _ datasource.DataSource = &AppStoreAppSettingsDataSource{}
var _ datasource.DataSourceWithConfigure = &AppStoreAppSettingsDataSource{}

func NewAppStoreAppSettingsDataSource() datasource.DataSource {
	return &AppStoreAppSettingsDataSource{}
}

// AppStoreAppSettingsDataSource is the data source implementation.
type AppStoreAppSettingsDataSource struct {
	client *sdk.Lumos
}

// AppStoreAppSettingsDataSourceModel describes the data model.
type AppStoreAppSettingsDataSourceModel struct {
	CustomRequestInstructions types.String                               `tfsdk:"custom_request_instructions"`
	ID                        types.String                               `tfsdk:"id"`
	Provisioning              *tfTypes.AddAppToAppStoreInputProvisioning `tfsdk:"provisioning"`
	RequestFlow               *tfTypes.AddAppToAppStoreInputRequestFlow  `tfsdk:"request_flow"`
}

// Metadata returns the data source type name.
func (r *AppStoreAppSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app_store_app_settings"
}

// Schema defines the schema for the data source.
func (r *AppStoreAppSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "AppStoreAppSettings DataSource",

		Attributes: map[string]schema.Attribute{
			"custom_request_instructions": schema.StringAttribute{
				Computed:    true,
				Description: `AppStore App instructions.`,
			},
			"id": schema.StringAttribute{
				Required: true,
			},
			"provisioning": schema.SingleNestedAttribute{
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
					"allow_multiple_permission_selection": schema.BoolAttribute{
						Computed:    true,
						Description: `Whether the app is configured to allow users to request multiple permissions in a single request`,
					},
					"custom_provisioning_instructions": schema.StringAttribute{
						Computed:    true,
						Description: `Only Available if manual steps is active. During the provisioning step, send a custom message to app admins explaining how to provision a user to the app. Markdown for links and text formatting is supported.`,
					},
					"groups_provisioning": schema.StringAttribute{
						Computed:    true,
						Description: `An enumeration. must be one of ["DIRECT_TO_USER", "GROUPS_AND_HIDDEN", "GROUPS_AND_VISIBLE"]`,
					},
					"manual_steps_needed": schema.BoolAttribute{
						Computed:    true,
						Description: `If enabled, Lumos will reach out to the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request.`,
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
				},
				Description: `Provisioning flow configuration to request access to app.`,
			},
			"request_flow": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"admins": schema.SingleNestedAttribute{
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
								Description: `Groups assigned as app admins.`,
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
								Description: `Users assigned as app admins.`,
							},
						},
						Description: `AppStore App admins assigned.`,
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
						Description: `During the approval step, send a custom message to requesters. Markdown for links and text formatting is supported.`,
					},
					"discoverability": schema.StringAttribute{
						Computed:    true,
						Description: `An enumeration. must be one of ["FULL", "LIMITED", "NONE"]`,
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
					"require_additional_approval": schema.BoolAttribute{
						Computed:    true,
						Description: `Only turn on when working with sensitive permissions to ensure a smooth employee experience.`,
					},
					"require_manager_approval": schema.BoolAttribute{
						Computed:    true,
						Description: `When a user makes an access request, require that their manager approves the request before moving on to additional approvals.`,
					},
				},
				Description: `Request flow configuration to request access to app.`,
			},
		},
	}
}

func (r *AppStoreAppSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *AppStoreAppSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *AppStoreAppSettingsDataSourceModel
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

	appID := data.ID.ValueString()
	request := operations.GetAppStoreAppSettingsRequest{
		AppID: appID,
	}
	res, err := r.client.AppStore.GetAppStoreAppSettings(ctx, request)
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
	if res.AppStoreAppSettingsOutput == nil {
		resp.Diagnostics.AddError("unexpected response from API. No response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedAppStoreAppSettingsOutput(res.AppStoreAppSettingsOutput)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}