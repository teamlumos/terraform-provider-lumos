package objectvalidators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	tfTypes "github.com/teamlumos/terraform-provider-lumos/internal/provider/types"
)

var _ validator.Object = ObjectRequestConfigInputValidatorValidator{}

// a set of possible TBA options
var timeBasedAccessOptions = []string{
	"Unlimited",
	"2 hours",
	"4 hours",
	"8 hours",
	"12 hours",
	"24 hours",
	"72 hours",
	"1 day",
	"3 days",
	"7 days",
	"14 days",
	"30 days",
	"90 days",
}

type ObjectRequestConfigInputValidatorValidator struct{}

// Description describes the validation in plain text formatting.
func (v ObjectRequestConfigInputValidatorValidator) Description(_ context.Context) string {
	return "Validates the request config input object to ensure override settings are not in conflict with populated fields"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v ObjectRequestConfigInputValidatorValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// Validate performs the validation.
func (v ObjectRequestConfigInputValidatorValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	requestConfig := tfTypes.RequestConfigInputCreate{}
	req.ConfigValue.As(ctx, &requestConfig, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

	// Check that none of the overrides are false or null with corresponding fields populated,
	// and that none of the overrides are true without values poplated
	allowedGroupsPopulated, allowedGroupsPopulatedProperly, errorMessage := AllowedGroupsIsPopulatedProperly(requestConfig.AllowedGroups)
	// Skip validation if the override value is unknown (e.g., from an unresolved variable)
	if !requestConfig.AllowedGroupsOverride.IsUnknown() {
		if requestConfig.AllowedGroupsOverride.ValueBool() == false {
			if allowedGroupsPopulated {
				AddAttributeErrorToResp(resp, req, "Invalid allowed groups", fmt.Sprintf("`allowed_groups_override` is %s but `allowed_groups` is populated", requestConfig.AllowedGroupsOverride.String()))
			}
		} else if !allowedGroupsPopulatedProperly {
			// allowed_groups_override is true
			AddAttributeErrorToResp(resp, req, "Invalid allowed groups", fmt.Sprintf("`allowed_groups_override` is true but `%s", errorMessage))
		}
	}

	if requestConfig.RequestApprovalConfig != nil {
		// Skip validation if the override value is unknown (e.g., from an unresolved variable)
		if !requestConfig.RequestApprovalConfig.CustomApprovalMessageOverride.IsUnknown() {
			if requestConfig.RequestApprovalConfig.CustomApprovalMessageOverride.ValueBool() == true {
				if requestConfig.RequestApprovalConfig.CustomApprovalMessage.IsNull() {
					resp.Diagnostics.AddAttributeError(
						req.Path,
						"Invalid custom approval message",
						fmt.Sprintf("`custom_approval_message_override` is true but `custom_approval_message` is not populated"),
					)
				}
			} else if !requestConfig.RequestApprovalConfig.CustomApprovalMessage.IsNull() {
				AddAttributeErrorToResp(resp,
					req,
					"Invalid custom approval message",
					fmt.Sprintf("`custom_approval_message_override` is %s but `custom_approval_message` is populated (%s)",
						requestConfig.RequestApprovalConfig.CustomApprovalMessageOverride.String(),
						requestConfig.RequestApprovalConfig.CustomApprovalMessage.String()),
				)
			}
		}

		// Skip validation if the override value is unknown (e.g., from an unresolved variable)
		if !requestConfig.RequestApprovalConfig.RequestApprovalConfigOverride.IsUnknown() {
			if requestConfig.RequestApprovalConfig.RequestApprovalConfigOverride.ValueBool() == true {
				if !ApproversIsPopulated(requestConfig.RequestApprovalConfig.Approvers) && requestConfig.RequestApprovalConfig.ManagerApproval.IsNull() {
					AddAttributeErrorToResp(resp,
						req,
						"Invalid request approval config",
						"`request_approval_config_override` is true but `approvers` and/or `manager_approval` are not populated",
					)
				}
			} else if ApproversIsPopulated(requestConfig.RequestApprovalConfig.Approvers) {
				AddAttributeErrorToResp(resp, req,
					"Invalid request approval config",
					fmt.Sprintf("`request_approval_config_override` is %s but `approvers` are populated", requestConfig.RequestApprovalConfig.RequestApprovalConfigOverride.String()),
				)
			} else if !requestConfig.RequestApprovalConfig.ManagerApproval.IsNull() {
				AddAttributeErrorToResp(resp, req,
					"Invalid request approval config",
					fmt.Sprintf("`request_approval_config_override` is %s but `manager_approval` is populated", requestConfig.RequestApprovalConfig.RequestApprovalConfigOverride.String()),
				)
			}
		}

		if !ApproversIsPopulated(requestConfig.RequestApprovalConfig.Approvers) && ApproversIsPopulated(requestConfig.RequestApprovalConfig.ApproversStage2) {
			AddAttributeErrorToResp(resp, req,
				"Invalid approvers",
				fmt.Sprintf("`approvers` is not populated but `approvers_stage_2`is populated (%s)", requestConfig.RequestApprovalConfig.ApproversStage2),
			)
		}
	}

	if requestConfig.RequestFulfillmentConfig != nil {
		// Skip validation if the override value is unknown (e.g., from an unresolved variable)
		if !requestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride.IsUnknown() {
			if requestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride.ValueBool() == true {
				if requestConfig.RequestFulfillmentConfig.TimeBasedAccess == nil || len(requestConfig.RequestFulfillmentConfig.TimeBasedAccess) == 0 {
					AddAttributeErrorToResp(resp, req,
						"Invalid time based access",
						fmt.Sprintf("`time_based_access_override` is true but `time_based_access` is not populated (%s)", requestConfig.RequestFulfillmentConfig.TimeBasedAccess),
					)
				}
				for i := 0; i < len(requestConfig.RequestFulfillmentConfig.TimeBasedAccess); i++ {
					opt := requestConfig.RequestFulfillmentConfig.TimeBasedAccess[i].ValueString()
					matchFound := false
					for j := 0; j < len(timeBasedAccessOptions); j++ {
						if opt == timeBasedAccessOptions[j] {
							matchFound = true
							break
						}
					}

					if !matchFound {
						AddAttributeErrorToResp(resp, req,
							"Invalid time based access option",
							fmt.Sprintf("`time_based_access` contains an invalid option (%s). All possible values (may or may not be configured for application): %s", opt, strings.Join(timeBasedAccessOptions, ", ")),
						)
					}
				}
			} else if requestConfig.RequestFulfillmentConfig.TimeBasedAccess != nil && len(requestConfig.RequestFulfillmentConfig.TimeBasedAccess) > 0 {
				AddAttributeErrorToResp(resp, req,
					"Invalid time based access",
					fmt.Sprintf("`time_based_access_override` is %s but `time_based_access` is populated (%s)", requestConfig.RequestFulfillmentConfig.TimeBasedAccessOverride.String(), requestConfig.RequestFulfillmentConfig.TimeBasedAccess),
				)
			}
		}
	}
}

func ApproversIsPopulated(approvers *tfTypes.AppApproversInput) bool {
	if approvers == nil {
		return false
	}

	if approvers.Users != nil && len(approvers.Users) > 0 {
		return true
	}

	return approvers.Groups != nil && len(approvers.Groups) > 0
}

// first bool indicates if its populated
// second bool indicate if its *properly* populated
func AllowedGroupsIsPopulatedProperly(allowedGroups *tfTypes.AllowedGroupsConfigInput) (bool, bool, string) {
	if allowedGroups == nil {
		return false, false, "`allowed_groups` is not populated"
	}

	if allowedGroups.Type.IsNull() {
		return true, false, "`allowed_groups` type is not populated"
	}

	if allowedGroups.Type.ValueString() != "ALL_GROUPS" && allowedGroups.Type.ValueString() != "SPECIFIED_GROUPS" {
		return true, false, fmt.Sprintf("`allowed_groups` type is invalid (%s)", allowedGroups.Type.ValueString())
	}

	if allowedGroups.Type.ValueString() == "ALL_GROUPS" {
		return true, true, ""
	}

	if allowedGroups.Groups != nil && len(allowedGroups.Groups) > 0 {
		if allowedGroups.Type.ValueString() != "SPECIFIED_GROUPS" {
			return true, false, "`allowed_groups` type is not populated properly for the groups specified"
		}
		return true, true, ""
	}

	return true, false, fmt.Sprintf("`allowed_groups` groups are not populated (%s)", allowedGroups.Groups)
}

func AddAttributeErrorToResp(resp *validator.ObjectResponse, req validator.ObjectRequest, attrName string, errMsg string) {
	resp.Diagnostics.AddAttributeError(req.Path, attrName, errMsg)
}

func RequestConfigInputValidator() validator.Object {
	return ObjectRequestConfigInputValidatorValidator{}
}
