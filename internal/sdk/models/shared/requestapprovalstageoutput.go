// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type RequestApprovalStageOutput struct {
	// The approvers of this stage.
	Approvers []ApproverOutput `json:"approvers,omitempty"`
}

func (o *RequestApprovalStageOutput) GetApprovers() []ApproverOutput {
	if o == nil {
		return nil
	}
	return o.Approvers
}
