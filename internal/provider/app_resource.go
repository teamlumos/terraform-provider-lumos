// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type appResourceModel struct {
	Name        types.String                 `tfsdk:"name"`
	Description types.String                 `tfsdk:"description"`
	Category    types.String                 `tfsdk:"category"`
	Settings    appStoreSettingResourceModel `tfsdk:"settings"`
}

type appStoreSettingResourceModel struct {
	CustomRequestInstructions types.String                `tfsdk:"custom_request_instructions"`
	RequestFlow               appStoreRequestFlowSetting  `tfsdk:"request_flow"`
	Provisioning              appStoreProvisioningSetting `tfsdk:"provisioning"`
	InAppStore                types.Bool                  `tfsdk:"in_app_store"`
}
