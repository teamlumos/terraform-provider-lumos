// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type LineItemInput struct {
	// The name of the line item as stored in Lumos
	Name string `json:"name"`
	// The type of purchase that this line item refers to
	Type string `json:"type"`
	// The number of units purchased for this line item
	Quantity int64 `json:"quantity"`
	// The per-unit cost of the line item
	UnitCost LineItemUnitCostInput `json:"unit_cost"`
}

func (o *LineItemInput) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *LineItemInput) GetType() string {
	if o == nil {
		return ""
	}
	return o.Type
}

func (o *LineItemInput) GetQuantity() int64 {
	if o == nil {
		return 0
	}
	return o.Quantity
}

func (o *LineItemInput) GetUnitCost() LineItemUnitCostInput {
	if o == nil {
		return LineItemUnitCostInput{}
	}
	return o.UnitCost
}
