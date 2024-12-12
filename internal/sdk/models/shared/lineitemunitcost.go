// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type LineItemUnitCost struct {
	// The per unit cost associated with this line item, amortized to the cost per month
	PerMonth Cost `json:"per_month"`
}

func (o *LineItemUnitCost) GetPerMonth() Cost {
	if o == nil {
		return Cost{}
	}
	return o.PerMonth
}