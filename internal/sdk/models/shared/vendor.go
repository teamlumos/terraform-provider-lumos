// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type Vendor struct {
	// A user friendly name for the vendor
	Name string `json:"name"`
}

func (o *Vendor) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}
