// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package types

import "github.com/hashicorp/terraform-plugin-framework/types"

type User struct {
	Email      types.String `tfsdk:"email"`
	FamilyName types.String `tfsdk:"family_name"`
	GivenName  types.String `tfsdk:"given_name"`
	ID         types.String `tfsdk:"id"`
	Status     types.String `tfsdk:"status"`
}
