// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type PageAccessRequest struct {
	Items []AccessRequest `json:"items"`
	Total *int64          `json:"total,omitempty"`
	Page  *int64          `json:"page,omitempty"`
	Size  *int64          `json:"size,omitempty"`
	Pages *int64          `json:"pages,omitempty"`
}

func (o *PageAccessRequest) GetItems() []AccessRequest {
	if o == nil {
		return []AccessRequest{}
	}
	return o.Items
}

func (o *PageAccessRequest) GetTotal() *int64 {
	if o == nil {
		return nil
	}
	return o.Total
}

func (o *PageAccessRequest) GetPage() *int64 {
	if o == nil {
		return nil
	}
	return o.Page
}

func (o *PageAccessRequest) GetSize() *int64 {
	if o == nil {
		return nil
	}
	return o.Size
}

func (o *PageAccessRequest) GetPages() *int64 {
	if o == nil {
		return nil
	}
	return o.Pages
}
