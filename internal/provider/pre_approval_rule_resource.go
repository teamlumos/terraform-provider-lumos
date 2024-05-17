package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/teamlumos/lumos-go-sdk/lumos"
	"github.com/teamlumos/terraform-provider-lumos/internal/generated/resource_pre_approval_rule"
)

var _ resource.Resource = (*preApprovalRuleResource)(nil)

func NewPreApprovalRuleResource() resource.Resource {
	return &preApprovalRuleResource{}
}

type preApprovalRuleResource struct {
	client *lumos.ClientWithResponses
}

func (r *preApprovalRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	lumos_client, ok := req.ProviderData.(*lumos.ClientWithResponses)
	if !ok {
		return
	}

	r.client = lumos_client
}

func (r *preApprovalRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pre_approval_rule"
}

func (r *preApprovalRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_pre_approval_rule.PreApprovalRuleResourceSchema(ctx)
}

func (r *preApprovalRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resource_pre_approval_rule.PreApprovalRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	/* Parse Terraform model resource_pre_approval_rule.PreApprovalRuleModel into lumos.PreApprovalRuleInput */
	apiInput, err := PreApprovalRuleTerraformToApi(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing payload", err.Error(),
		)
		return

	}

	/* Make API call */
	lumosApiResponse, err := r.client.CreatePreApprovalRuleAppstorePreApprovalRulesPostWithResponse(ctx, *apiInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing payload", err.Error(),
		)
		return
	}

	if lumosApiResponse.HTTPResponse.StatusCode != 201 {
		errorMessage := fmt.Sprintf("Expected status code 201, got %d. %s", lumosApiResponse.HTTPResponse.StatusCode, lumosApiResponse.Body)
		resp.Diagnostics.AddError(
			"Unsuccessful API call", errorMessage,
		)
		return
	}

	/* Parse API response lumos.CreatePreApprovalRuleAppstorePreApprovalRulesPostResponse into resource_pre_approval_rule.PreApprovalRuleModel */
	*data = PreApprovalRuleApiToTerraform(ctx, *lumosApiResponse.JSON201, *data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *preApprovalRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resource_pre_approval_rule.PreApprovalRuleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	lumosApiResponse, err := r.client.GetAppstorePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIdGetWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pre approval rule",
			"Could not get pre approval rule "+data.Id.ValueString()+" unexpected error:"+err.Error(),
		)
		return
	}
	*data = PreApprovalRuleApiToTerraform(ctx, *lumosApiResponse.JSON200, *data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *preApprovalRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *resource_pre_approval_rule.PreApprovalRuleModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiInput, err := PreApprovalRuleTerraformToApi(ctx, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing payload", err.Error(),
		)
		return

	}
	apiUpdateInput := lumos.PreApprovalRuleUpdateInput{
		Justification:          apiInput.Justification,
		PreapprovalWebhooks:    apiInput.PreapprovalWebhooks,
		PreapprovedGroups:      apiInput.PreapprovedGroups,
		PreapprovedPermissions: apiInput.PreapprovedPermissions,
		TimeBasedAccess:        apiInput.TimeBasedAccess,
	}

	lumosApiResponse, err := r.client.UpdatePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIdPatchWithResponse(ctx, data.Id.ValueString(), apiUpdateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing payload", err.Error(),
		)
		return
	}

	if lumosApiResponse.HTTPResponse.StatusCode != 200 {
		errorMessage := fmt.Sprintf("Expected status code 200, got %d. %s", lumosApiResponse.HTTPResponse.StatusCode, lumosApiResponse.Body)
		resp.Diagnostics.AddError(
			"Unsuccessful API call", errorMessage,
		)
		return
	}

	*data = PreApprovalRuleApiToTerraform(ctx, *lumosApiResponse.JSON200, *data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *preApprovalRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_pre_approval_rule.PreApprovalRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteAppstorePermissionAppstoreRequestablePermissionsPermissionIdDelete(ctx, data.Id.ValueString())
	if err != nil {
		return
	}
}
