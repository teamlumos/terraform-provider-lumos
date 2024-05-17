package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/teamlumos/lumos-go-sdk/lumos"
	"github.com/teamlumos/terraform-provider-lumos/internal/generated/datasource_pre_approval_rule"
)

var _ datasource.DataSource = (*preApprovalRuleDataSource)(nil)

func NewPreApprovalRuleDataSource() datasource.DataSource {
	return &preApprovalRuleDataSource{}
}

type preApprovalRuleDataSource struct {
	client *lumos.ClientWithResponses
}

func (d *preApprovalRuleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	lumos_client, ok := req.ProviderData.(*lumos.ClientWithResponses)
	if !ok {
		return
	}

	d.client = lumos_client
}

func (d *preApprovalRuleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pre_approval_rule"
}

func (d *preApprovalRuleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_pre_approval_rule.PreApprovalRuleDataSourceSchema(ctx)
}

func (d *preApprovalRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *datasource_pre_approval_rule.PreApprovalRuleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	lumosApiResponse, err := d.client.GetAppstorePreApprovalRuleAppstorePreApprovalRulesPreApprovalRuleIdGetWithResponse(ctx, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting pre approval rule",
			"Could not get pre approval rule "+data.Id.ValueString()+" unexpected error:"+err.Error(),
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

	data.Id = types.StringPointerValue(lumosApiResponse.JSON200.Id)
	data.AppClassId = types.StringPointerValue(lumosApiResponse.JSON200.AppClassId)
	data.AppId = types.StringPointerValue(lumosApiResponse.JSON200.AppId)
	data.AppInstanceId = types.StringPointerValue(lumosApiResponse.JSON200.AppInstanceId)
	data.Justification = types.StringValue("bullshit") //types.StringPointerValue(&lumosApiResponse.JSON200.Justification)
	var timeBasedAccess = []string{}
	for i, timeBasedAccessOption := range *lumosApiResponse.JSON200.TimeBasedAccess {
		timeBasedAccess[i] = string(timeBasedAccessOption)
	}
	// data.TimeBasedAccess, _ = types.ListValueFrom(ctx, types.StringType, timeBasedAccess)

	/*
		In the resources case, we know that the length of the API response is equal to the length of the attribute in the Terraform model.

		However in this case, the length of the Terraform model for groups, permissions, webhooks is 0, but the API response may be longer.

		I think it is critical for us to run ElementsAs so that ListValueFrom can properly pick up the types
		But it's difficult because that method is on the Terraform model, which we know will be empty.
		Almost chicken and the egg.
	*/
	var apiResponseLength = len(*lumosApiResponse.JSON200.PreapprovedGroups)
	var preApprovedGroups = make([]datasource_pre_approval_rule.PreapprovedGroupsValue, apiResponseLength)
	for i, group := range *lumosApiResponse.JSON200.PreapprovedGroups {
		preApprovedGroups[i].Id = types.StringPointerValue(group.Id)
	}

	/*
		var preApprovedGroups []datasource_pre_approval_rule.PreapprovedGroupsValue
		data.PreapprovedGroups.ElementsAs(ctx, &preApprovedGroups, false)
		for _, group := range *lumosApiResponse.JSON200.PreapprovedGroups {
			value := datasource_pre_approval_rule.PreapprovedGroupsValue{
				Id:                    types.StringPointerValue(group.Id),
				AppId:                 types.StringPointerValue(group.AppId),
				IntegrationSpecificId: types.StringPointerValue(group.IntegrationSpecificId),
				Name:                  types.StringPointerValue(group.Name),
				Description:           types.StringPointerValue(group.Description),
				Lifecycle:             types.StringPointerValue((*string)(group.Lifecycle)),
				SourceAppId:           types.StringPointerValue(group.SourceAppId),
			}
			preApprovedGroups = append(preApprovedGroups, value)
		}
		if len(preApprovedGroups) > 0 {
			// var d diag.Diagnostics
			res, _ := types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), preApprovedGroups)
			// data.PreapprovedGroups, d = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), preApprovedGroups)
			other_res, _ := types.ListValueFrom(ctx, resource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), preApprovedGroups)
			fmt.Printf("RESULT: %#+v\n", res)
			fmt.Printf("OTHER RESULT: %#+v\n", other_res)
		} else {
			var d diag.Diagnostics
			data.PreapprovedGroups, d = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), make([]datasource_pre_approval_rule.PreapprovedGroupsValue, 0))
			fmt.Printf("2DIAG: %#+v\n", d)
		}
		fmt.Printf("2preApprovedGroups: %#+v\n", preApprovedGroups)
		fmt.Printf("len2preApprovedGroups: %#+v\n", len(preApprovedGroups))
		fmt.Printf("333preapproved: %#+v\n", data.PreapprovedGroups)
	*/
	/*
		for i, group := range *lumosApiResponse.JSON200.PreapprovedGroups {
			preApprovedGroups[i].Id = types.StringPointerValue(group.Id)
		}
	*/
	/*
		if len(preApprovedGroups) > 0 {
			data.PreapprovedGroups, _ = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), preApprovedGroups)
		} else {
			data.PreapprovedGroups, _ = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedGroupsValue{}.Type(ctx), make([]datasource_pre_approval_rule.PreapprovedGroupsValue, 0))
		}

			var preApprovedPermissions []datasource_pre_approval_rule.PreapprovedPermissionsValue
			data.PreapprovedPermissions.ElementsAs(ctx, &preApprovedPermissions, false)
			for i, permission := range *lumosApiResponse.JSON200.PreapprovedPermissions {
				preApprovedPermissions[i].Id = types.StringPointerValue(permission.Id)
				preApprovedPermissions[i].Label = types.StringValue(permission.Label)
				preApprovedPermissions[i].PreapprovedPermissionsType = types.StringPointerValue((*string)(permission.Type))
				preApprovedPermissions[i].AppId = types.StringValue(permission.AppId)
				preApprovedPermissions[i].AppClassId = types.StringValue(permission.AppClassId)
				preApprovedPermissions[i].AppInstanceId = types.StringValue(permission.AppInstanceId)
			}
			if len(preApprovedPermissions) > 0 {
				data.PreapprovedPermissions, _ = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedPermissionsValue{}.Type(ctx), preApprovedPermissions)
			} else {
				data.PreapprovedPermissions, _ = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovedPermissionsValue{}.Type(ctx), make([]datasource_pre_approval_rule.PreapprovedPermissionsValue, 0))

			}

			var preApprovalWebhooks []datasource_pre_approval_rule.PreapprovalWebhooksValue
			data.PreapprovalWebhooks.ElementsAs(ctx, &preApprovalWebhooks, false)
			for i, webhook := range *lumosApiResponse.JSON200.PreapprovalWebhooks {
				preApprovalWebhooks[i].Description = types.StringPointerValue(webhook.Description)
				preApprovalWebhooks[i].HookType = types.StringValue(string(webhook.HookType))
				preApprovalWebhooks[i].Id = types.StringValue(string(webhook.Id))
				preApprovalWebhooks[i].Name = types.StringValue(string(webhook.Name))
			}
			if len(preApprovalWebhooks) > 0 {
				data.PreapprovalWebhooks, _ = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovalWebhooksValue{}.Type(ctx), preApprovalWebhooks)
			} else {
				data.PreapprovalWebhooks, _ = types.ListValueFrom(ctx, datasource_pre_approval_rule.PreapprovalWebhooksValue{}.Type(ctx), make([]datasource_pre_approval_rule.PreapprovalWebhooksValue, 0))
			}
	*/

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
