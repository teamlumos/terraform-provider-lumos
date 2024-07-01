// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type PagePreApprovalRuleOutput struct {
	Items []PreApprovalRuleOutput `json:"items"`
	Page  *int64                  `json:"page,omitempty"`
	Pages *int64                  `json:"pages,omitempty"`
	Size  *int64                  `json:"size,omitempty"`
	Total *int64                  `json:"total,omitempty"`
}

func (o *PagePreApprovalRuleOutput) GetItems() []PreApprovalRuleOutput {
	if o == nil {
		return []PreApprovalRuleOutput{}
	}
	return o.Items
}

func (o *PagePreApprovalRuleOutput) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *PagePreApprovalRuleOutput) GetPages() *int64 {
	if o == nil {
		return nil
	}
	return o.Pages
}

func (o *PagePreApprovalRuleOutput) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

func (o *PagePreApprovalRuleOutput) GetTotal() *int64 {
	if o == nil {
		return nil
	}
	return o.Total
}
