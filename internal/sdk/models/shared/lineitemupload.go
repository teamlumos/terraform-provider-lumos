// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

// LineItemUploadUnitCost - The unit cost of this line item
type LineItemUploadUnitCost struct {
	// The currency in which this cost is stored
	Currency *string `json:"currency,omitempty"`
	// The quantity of the cost in terms of the specified currency
	Value int64 `json:"value"`
}

func (o *LineItemUploadUnitCost) GetCurrency() *string {
	if o == nil {
		return nil
	}
	return o.Currency
}

func (o *LineItemUploadUnitCost) GetValue() int64 {
	if o == nil {
		return 0
	}
	return o.Value
}

type LineItemUpload struct {
	// The name of the line item as stored in Lumos
	Name string `json:"name"`
	// The number of units purchased for this line item
	Quantity int64 `json:"quantity"`
	// The unit cost of this line item
	UnitCost LineItemUploadUnitCost `json:"unit_cost"`
}

func (o *LineItemUpload) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

func (o *LineItemUpload) GetQuantity() int64 {
	if o == nil {
		return 0
	}
	return o.Quantity
}

func (o *LineItemUpload) GetUnitCost() LineItemUploadUnitCost {
	if o == nil {
		return LineItemUploadUnitCost{}
	}
	return o.UnitCost
}