// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type Data struct {
}

// AirbasePurchaseRequestEvent - Model of the input we'll get from Airbase's webhook call.
// All events seem to follow this same base structure, it is just `data` that is different
// for each. We are using `dict` here to be permissive to other webhook events (so we can warn if
// the customer misconfigures for a different webhook), and we actually only use 2 fields from it.
type AirbasePurchaseRequestEvent struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Type        string `json:"type"`
	CreatedDate string `json:"created_date"`
	Data        Data   `json:"data"`
}

func (o *AirbasePurchaseRequestEvent) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *AirbasePurchaseRequestEvent) GetObject() string {
	if o == nil {
		return ""
	}
	return o.Object
}

func (o *AirbasePurchaseRequestEvent) GetType() string {
	if o == nil {
		return ""
	}
	return o.Type
}

func (o *AirbasePurchaseRequestEvent) GetCreatedDate() string {
	if o == nil {
		return ""
	}
	return o.CreatedDate
}

func (o *AirbasePurchaseRequestEvent) GetData() Data {
	if o == nil {
		return Data{}
	}
	return o.Data
}
