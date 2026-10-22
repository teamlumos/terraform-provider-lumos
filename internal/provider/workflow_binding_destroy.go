package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Lumos can't empty these workflow slots, so destroying one of these
// bindings, or moving it to another app or entitlement, would leave the
// workflow bound where it was while Terraform reported it gone. The plan
// fails instead and says how to change the slot or stop managing it. A
// removed block with destroy = false never plans a destroy, so it isn't
// affected. lumos_entitlement_approval_workflow isn't listed: destroying it
// turns the override off.

var (
	_ resource.ResourceWithModifyPlan = &AppApprovalWorkflowResource{}
	_ resource.ResourceWithModifyPlan = &AppProvisioningWorkflowResource{}
	_ resource.ResourceWithModifyPlan = &AppDeprovisioningWorkflowResource{}
	_ resource.ResourceWithModifyPlan = &EntitlementProvisioningWorkflowResource{}
	_ resource.ResourceWithModifyPlan = &EntitlementDeprovisioningWorkflowResource{}
)

func (r *AppApprovalWorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	refuseUnbinding(ctx, req, resp, "lumos_app_approval_workflow", "app_id", "app")
}

func (r *AppProvisioningWorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	refuseUnbinding(ctx, req, resp, "lumos_app_provisioning_workflow", "app_id", "app")
}

func (r *AppDeprovisioningWorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	refuseUnbinding(ctx, req, resp, "lumos_app_deprovisioning_workflow", "app_id", "app")
}

func (r *EntitlementProvisioningWorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	refuseUnbinding(ctx, req, resp, "lumos_entitlement_provisioning_workflow", "entitlement_id", "entitlement")
}

func (r *EntitlementDeprovisioningWorkflowResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	refuseUnbinding(ctx, req, resp, "lumos_entitlement_deprovisioning_workflow", "entitlement_id", "entitlement")
}

// refuseUnbinding fails a plan that destroys the binding or moves it to a
// different owner. Creates and workflow_id changes plan as usual.
func refuseUnbinding(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, typeName, ownerAttr, ownerKind string) {
	if req.State.Raw.IsNull() {
		return
	}
	var ownerID, workflowID types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(ownerAttr), &ownerID)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("workflow_id"), &workflowID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	help := fmt.Sprintf("To change it, set workflow_id to a different workflow. To stop managing it "+
		"without changing Lumos, replace the resource block with a removed block whose lifecycle sets "+
		"destroy = false.\n\nThe %s keeps workflow %s.", ownerKind, workflowID.ValueString())

	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddError(
			"Workflow binding can't be removed",
			fmt.Sprintf("Destroying this %s wouldn't unbind workflow %s from %s %s: Lumos can't empty "+
				"this workflow slot. %s", typeName, workflowID.ValueString(), ownerKind, ownerID.ValueString(), help),
		)
		return
	}
	var plannedOwner types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(ownerAttr), &plannedOwner)...)
	if resp.Diagnostics.HasError() || plannedOwner.IsUnknown() || plannedOwner.Equal(ownerID) {
		return
	}
	resp.Diagnostics.AddAttributeError(
		path.Root(ownerAttr),
		"Workflow binding can't be moved",
		fmt.Sprintf("Changing %s on this %s would replace it, but destroying the old binding wouldn't "+
			"unbind workflow %s from %s %s: Lumos can't empty this workflow slot. To bind a workflow to %s "+
			"%s, add a separate resource for it. %s", ownerAttr, typeName, workflowID.ValueString(),
			ownerKind, ownerID.ValueString(), ownerKind, plannedOwner.ValueString(), help),
	)
}
