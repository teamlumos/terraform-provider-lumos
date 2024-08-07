// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type Links struct {
	First *string `json:"first,omitempty"`
	Last  *string `json:"last,omitempty"`
	Self  *string `json:"self,omitempty"`
	Next  *string `json:"next,omitempty"`
	Prev  *string `json:"prev,omitempty"`
}

func (o *Links) GetFirst() *string {
	if o == nil {
		return nil
	}
	return o.First
}

func (o *Links) GetLast() *string {
	if o == nil {
		return nil
	}
	return o.Last
}

func (o *Links) GetSelf() *string {
	if o == nil {
		return nil
	}
	return o.Self
}

func (o *Links) GetNext() *string {
	if o == nil {
		return nil
	}
	return o.Next
}

func (o *Links) GetPrev() *string {
	if o == nil {
		return nil
	}
	return o.Prev
}
