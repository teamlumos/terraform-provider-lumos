// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &requestablePermissionResource{}
var _ resource.ResourceWithImportState = &requestablePermissionResource{}

func NewRequestablePermissionResource() resource.Resource {
	return &requestablePermissionResource{}
}

type requestablePermissionResource struct {
	client *LumosAPIClient
}

type ProvisioningGroupModel struct {
	Id                    types.String `tfsdk:"id"`
	AppId                 types.String `tfsdk:"app_id"`
	IntegrationSpecificId types.String `tfsdk:"integration_specific_id"`
}

type requestablePermissionResourceModel struct {
	Id                             types.String           `tfsdk:"id"`
	AppId                          types.String           `tfsdk:"app_id"`
	Label                          types.String           `tfsdk:"label"`
	PermissionType                 types.String           `tfsdk:"permission_type"`
	VisibleInAppStore              types.Bool             `tfsdk:"visible_in_appstore"`
	AllowedGroups                  types.Set              `tfsdk:"allowed_groups"`
	ManagerApprovalRequired        types.Bool             `tfsdk:"manager_approval_required"`
	ApproverGroupsStage1           types.Set              `tfsdk:"approver_groups_stage_1"`
	ApproverUsersStage1            types.Set              `tfsdk:"approver_users_stage_1"`
	ApproverGroupsStage2           types.Set              `tfsdk:"approver_groups_stage_2"`
	ApproverUsersStage2            types.Set              `tfsdk:"approver_users_stage_2"`
	RequestApprovalConfigOverride  types.Bool             `tfsdk:"request_approval_config_override"`
	RequireAdditionalApproval      types.Bool             `tfsdk:"require_additional_approval"`
	AccessRemovalInlineWebhook     types.String           `tfsdk:"access_removal_inline_webhook"`
	RequestValidationInlineWebhook types.String           `tfsdk:"request_validation_inline_webhook"`
	ProvisioningInlineWebhook      types.String           `tfsdk:"provisioning_inline_webhook"`
	ManualStepsNeeded              types.Bool             `tfsdk:"manual_steps_needed"`
	ManualInstructions             types.String           `tfsdk:"manual_instructions"`
	TimeBasedAccessOptions         types.Set              `tfsdk:"time_based_access_options"`
	TimeBasedAccessOverride        types.Bool             `tfsdk:"time_based_access_override"`
	ProvisioningGroup              ProvisioningGroupModel `tfsdk:"provisioning_group"`
	LastUpdated                    types.String           `tfsdk:"last_updated"`
}

func (r *requestablePermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_requestable_permission"
}

func (r *requestablePermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Requestable Permission",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The ID of this requestable permission.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"app_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the app associated with this requestable permission",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				Required:    true,
				Description: "The label of this requestable permission.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"permission_type": schema.StringAttribute{
				Computed:    true,
				Description: "Define if its a 'SYNCED' permission from idp with a real association with the app or a 'NATIVE' [manual] permission.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"visible_in_appstore": schema.BoolAttribute{
				Required:    true,
				Description: "Defines if the permission will be visible in the AppStore to request access or not.",
			},
			"manager_approval_required": schema.BoolAttribute{
				Required:    true,
				Description: "When a user makes an access request, require that their manager approves the request before moving on to additional approvals.",
			},
			"allowed_groups": schema.SetAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "Define if its a 'SYNCED' permission from idp with a real association with the app or a 'NATIVE' [manual] permission.",
			},
			"approver_groups_stage_1": schema.SetAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "Group IDs assigned as approvers.",
			},
			"approver_users_stage_1": schema.SetAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "User IDs assigned as approvers.",
			},
			"approver_groups_stage_2": schema.SetAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "Group IDs assigned as approvers in stage 2",
			},
			"approver_users_stage_2": schema.SetAttribute{
				Required:    false,
				Optional:    true,
				ElementType: types.StringType,
				Description: "User IDs assigned as approvers in stage 2",
			},
			"require_additional_approval": schema.BoolAttribute{
				Required:    true,
				Description: "Only turn on when working with sensitive permissions to ensure a smooth employee experience.",
			},
			"request_approval_config_override": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "TODO",
			},
			"request_validation_inline_webhook": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "The ID of request validation webhook can be optionally associated with this app.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provisioning_inline_webhook": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "The ID of provisioning webhook optionally associated with this app.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_removal_inline_webhook": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: " The ID of inactivity workflow can be optionally associated with this app.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"time_based_access_options": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "Access length options available for the users to request an app for a selected duration. After expiry, Lumos will automatically remove user's access. If empty default will be Unlimited.",
			},
			"time_based_access_override": schema.BoolAttribute{
				Required:            true,
				MarkdownDescription: "TODO",
			},
			"manual_steps_needed": schema.BoolAttribute{
				Required:    true,
				Description: "If enabled, Lumos will reach out to the App Admin after initial access is granted to perform additional manual steps. Note that if this option is enabled, this action must be confirmed by the App Admin in order to resolve the request.",
			},
			"manual_instructions": schema.StringAttribute{
				Required:    false,
				Optional:    true,
				Description: "Only Available if manual steps is active. During the provisioning step, send a custom message to app admins explaining how to provision a user to the permission. Markdown for links and text formatting is supported.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"provisioning_group": schema.SingleNestedAttribute{
				Required:            false,
				Optional:            true,
				MarkdownDescription: "The provisioning group ID optionally associated with this config",
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Required:            false,
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The ID of the provisioning group",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"app_id": schema.StringAttribute{
						Required:            false,
						Optional:            true,
						MarkdownDescription: "The ID of the app associated with this provisioning group",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"integration_specific_id": schema.StringAttribute{
						Required:            false,
						Optional:            true,
						MarkdownDescription: "The ID of the integration specific to the provisioning group",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *requestablePermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func validateOktaGroupId(groupId string) error {
	pattern := `^00g[a-zA-Z0-9]{17}$`

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("Invalid regex pattern")
	}

	if !regex.MatchString(groupId) {
		return fmt.Errorf("'%s' is not a valid Okta Group ID.", groupId)
	}

	return nil
}

func validateProvisioningGroup(provisioningGroup ProvisioningGroupModel) error {
	if provisioningGroup.Id.IsUnknown() || provisioningGroup.Id.IsNull() {
		if provisioningGroup.AppId.IsUnknown() || provisioningGroup.AppId.IsNull() || provisioningGroup.IntegrationSpecificId.IsUnknown() || provisioningGroup.IntegrationSpecificId.IsNull() {
			return fmt.Errorf("either 'id' must be specified or both 'app_id' and 'integration_specific_id' must be provided")
		}
	} else {
		if !provisioningGroup.AppId.IsNull() || !provisioningGroup.IntegrationSpecificId.IsNull() {
			return fmt.Errorf("'id' cannot be provided with 'app_id' or 'integration_specific_id'")
		}
	}

	// For now we expect all "IntegrationSpecificId"s to take the form of Okta group IDs
	if !provisioningGroup.IntegrationSpecificId.IsNull() {
		if err := validateOktaGroupId(provisioningGroup.IntegrationSpecificId.ValueString()); err != nil {
			return err
		}
	}

	return nil
}

func (r *requestablePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan requestablePermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "Error at the start.")
		return
	}

	if err := validateProvisioningGroup(plan.ProvisioningGroup); err != nil {
		resp.Diagnostics.AddError("Invalid Provisioning Group Configuration", err.Error())
		return
	}

	permission, err := r.client.createPermission(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating permission",
			"Could not create permission, unexpected error: "+err.Error(),
		)
		return
	}
	plan = setRequestablePermissionResourceModelFromLumosAPIRequestablePermission(ctx, plan, *permission)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *requestablePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *requestablePermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.getPermission(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting permission",
			"Could not get permission "+data.Id.ValueString()+" unexpected error:"+err.Error(),
		)
		return
	}

	*data = setRequestablePermissionResourceModelFromLumosAPIRequestablePermission(ctx, *data, *permission)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestablePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *requestablePermissionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.updatePermission(data.Id.ValueString(), *data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating permission",
			"Could not update permission "+data.Id.ValueString()+" unexpected error:"+err.Error(),
		)
		return
	}
	*data = setRequestablePermissionResourceModelFromLumosAPIRequestablePermission(ctx, *data, *permission)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *requestablePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data requestablePermissionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "Error at the start.")
		return
	}

	r.client.deletePermission(data.Id.ValueString())
}

func (r *requestablePermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
