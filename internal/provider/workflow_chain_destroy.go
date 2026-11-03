package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Lumos never deletes an approval chain, so destroying or replacing a
// lumos_workflow of type approval-chain only removes the old chain from
// state (see the approvalChainKept SDK hook). The plan warns, so that's not
// a surprise after the apply.

var _ resource.ResourceWithModifyPlan = &WorkflowResource{}

func (r *WorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() {
		return
	}
	var workflowType, id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("type"), &workflowType)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() || workflowType.ValueString() != "approval-chain" {
		return
	}
	action := "Replacing"
	if req.Plan.Raw.IsNull() {
		action = "Destroying"
	} else if len(resp.RequiresReplace) == 0 {
		return
	}
	resp.Diagnostics.AddWarning(
		"Lumos keeps approval chains",
		fmt.Sprintf("%s this lumos_workflow removes approval chain %s from Terraform state, but Lumos "+
			"doesn't delete approval chains, so the chain stays in Lumos.", action, id.ValueString()),
	)
}
