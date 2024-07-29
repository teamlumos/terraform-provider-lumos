// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	speakeasy_stringplanmodifier "github.com/teamlumos/terraform-provider-lumos/internal/planmodifiers/stringplanmodifier"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk"
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/operations"
	speakeasy_stringvalidators "github.com/teamlumos/terraform-provider-lumos/internal/validators/stringvalidators"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PreApprovalRuleResource{}
var _ resource.ResourceWithImportState = &PreApprovalRuleResource{}

func NewPreApprovalRuleResource() resource.Resource {
	return &PreApprovalRuleResource{}
}

// PreApprovalRuleResource defines the resource implementation.
type PreApprovalRuleResource struct {
	client *sdk.Lumos
}

// PreApprovalRuleResourceModel describes the resource data model.
type PreApprovalRuleResourceModel struct {
	AppClassID             types.String                                              `tfsdk:"app_class_id"`
	AppID                  types.String                                              `tfsdk:"app_id"`
	AppInstanceID          types.String                                              `tfsdk:"app_instance_id"`
	ID                     types.String                                              `tfsdk:"id"`
	Justification          types.String                                              `tfsdk:"justification"`
	PreapprovalWebhooks    []tfTypes.AddAppToAppStoreInputAccessRemovalInlineWebhook `tfsdk:"preapproval_webhooks"`
	PreapprovedGroups      []tfTypes.Group                                           `tfsdk:"preapproved_groups"`
	PreapprovedPermissions []tfTypes.RequestablePermissionBase                       `tfsdk:"preapproved_permissions"`
	TimeBasedAccess        []types.String                                            `tfsdk:"time_based_access"`
}

func (r *PreApprovalRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pre_approval_rule"
}

func (r *PreApprovalRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "PreApprovalRule Resource",
		Attributes: map[string]schema.Attribute{
			"app_class_id": schema.StringAttribute{
				Computed:    true,
				Description: `The ID of the service associated with this pre-approval rule.`,
			},
			"app_id": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					speakeasy_stringplanmodifier.SuppressDiff(speakeasy_stringplanmodifier.ExplicitSuppress),
				},
				Required:    true,
				Description: `The ID of the app associated with this pre-approval rule. Requires replacement if changed. `,
			},
			"app_instance_id": schema.StringAttribute{
				Computed:    true,
				Description: `Optionally, an app has an identifer associated with it's particular instance.`,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					speakeasy_stringplanmodifier.SuppressDiff(speakeasy_stringplanmodifier.ExplicitSuppress),
				},
				Description: `The ID of this preapproval rule.`,
			},
			"justification": schema.StringAttribute{
				Required:    true,
				Description: `The justification of this preapproval rule.`,
			},
			"preapproval_webhooks": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Computed:    true,
							Description: `The description of this inline webhook.`,
						},
						"hook_type": schema.StringAttribute{
							Computed:    true,
							Description: `An enumeration. must be one of ["PRE_APPROVAL", "PROVISION", "DEPROVISION", "REQUEST_VALIDATION", "SIEM"]`,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"PRE_APPROVAL",
									"PROVISION",
									"DEPROVISION",
									"REQUEST_VALIDATION",
									"SIEM",
								),
							},
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: `The ID of this inline webhook. Not Null`,
							Validators: []validator.String{
								speakeasy_stringvalidators.NotNull(),
							},
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: `The name of this inline webhook.`,
						},
					},
				},
				Description: `The preapproval webhooks of this preapproval rule.`,
			},
			"preapproved_groups": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"app_id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: `The ID of the app that sources this group.`,
						},
						"description": schema.StringAttribute{
							Computed:    true,
							Description: `The description of this group.`,
						},
						"group_lifecycle": schema.StringAttribute{
							Computed:    true,
							Description: `The lifecycle of this group. must be one of ["SYNCED", "NATIVE"]`,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"SYNCED",
									"NATIVE",
								),
							},
						},
						"id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: `The ID of this group.`,
						},
						"integration_specific_id": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
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
				Description: `The preapproved groups of this preapproval rule.`,
			},
			"preapproved_permissions": schema.ListNestedAttribute{
				Computed: true,
				Optional: true,
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
							Optional:    true,
							Description: `The ID of this requestable permission.`,
						},
						"label": schema.StringAttribute{
							Computed:    true,
							Description: `The label of this requestable permission.`,
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: `An enumeration. must be one of ["SYNCED", "NATIVE"]`,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"SYNCED",
									"NATIVE",
								),
							},
						},
					},
				},
				Description: `The preapproved permissions of this preapproval rule.`,
			},
			"time_based_access": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				Description: `Preapproval rule time access length,`,
			},
		},
	}
}

func (r *PreApprovalRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.Lumos)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sdk.Lumos, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PreApprovalRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *PreApprovalRuleResourceModel
	var plan types.Object

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(plan.As(ctx, &data, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := *data.ToSharedPreApprovalRuleInput()
	res, err := r.client.AppStore.CreatePreApprovalRuleAppstorePreApprovalRulesPost(ctx, request)
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
	if res.StatusCode != 201 {
		resp.Diagnostics.AddError(fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", res.StatusCode), debugResponse(res.RawResponse))
		return
	}
	if !(res.PreApprovalRuleOutput != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedPreApprovalRuleOutput(res.PreApprovalRuleOutput)
	refreshPlan(ctx, plan, &data, resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PreApprovalRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *PreApprovalRuleResourceModel
	var item types.Object

	resp.Diagnostics.Append(req.State.Get(ctx, &item)...)
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

	// Not Implemented; we rely entirely on CREATE API request response

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PreApprovalRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *PreApprovalRuleResourceModel
	var plan types.Object

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	merge(ctx, req, resp, &data)
	if resp.Diagnostics.HasError() {
		return
	}

	var id string
	id = data.ID.ValueString()

	preApprovalRuleUpdateInput := *data.ToSharedPreApprovalRuleUpdateInput()
	request := operations.UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatchRequest{
		ID:                         id,
		PreApprovalRuleUpdateInput: preApprovalRuleUpdateInput,
	}
	res, err := r.client.AppStore.UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDPatch(ctx, request)
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
	if res.StatusCode != 200 {
		resp.Diagnostics.AddError(fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", res.StatusCode), debugResponse(res.RawResponse))
		return
	}
	if !(res.PreApprovalRuleOutput != nil) {
		resp.Diagnostics.AddError("unexpected response from API. Got an unexpected response body", debugResponse(res.RawResponse))
		return
	}
	data.RefreshFromSharedPreApprovalRuleOutput(res.PreApprovalRuleOutput)
	refreshPlan(ctx, plan, &data, resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PreApprovalRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *PreApprovalRuleResourceModel
	var item types.Object

	resp.Diagnostics.Append(req.State.Get(ctx, &item)...)
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

	var id string
	id = data.ID.ValueString()

	request := operations.DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDeleteRequest{
		ID: id,
	}
	res, err := r.client.AppStore.DeletePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIDDelete(ctx, request)
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
	if res.StatusCode != 204 {
		resp.Diagnostics.AddError(fmt.Sprintf("unexpected response from API. Got an unexpected response code %v", res.StatusCode), debugResponse(res.RawResponse))
		return
	}

}

func (r *PreApprovalRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.AddError("Not Implemented", "No available import state operation is available for resource pre_approval_rule.")
}
