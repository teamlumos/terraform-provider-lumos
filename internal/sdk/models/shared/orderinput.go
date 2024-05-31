// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

import (
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/internal/utils"
	"github.com/speakeasy/terraform-provider-lumos/internal/sdk/types"
)

// OrderInputVendor - Information about the vendor associated with this contract. Currently, we will only support a name field.
type OrderInputVendor struct {
	// A user friendly name for the vendor
	Name string `json:"name"`
}

func (o *OrderInputVendor) GetName() string {
	if o == nil {
		return ""
	}
	return o.Name
}

type OrderInput struct {
	// A unique ID for the contract being uploaded. This can be an ID from an external system like Ironclad, an internal ID such as a PO number, or simply the name of the vendor + date of the contract.
	UniqueIdentifier string `json:"unique_identifier"`
	// Information about the vendor associated with this contract. Currently, we will only support a name field.
	Vendor OrderInputVendor `json:"vendor"`
	// The start date of the attached contract
	StartDate types.Date `json:"start_date"`
	// The end date of the attached contract
	EndDate types.Date `json:"end_date"`
	// The date by which a vendor must be notified before terminating the order
	OptOutDate *types.Date `json:"opt_out_date,omitempty"`
	// Whether or not the contract auto-renews
	AutoRenewal bool `json:"auto_renewal"`
	// The currency in which the contract is being paid in ISO 4217 format
	Currency *string `json:"currency,omitempty"`
	// The list of currently active line items for this Vendor Agreement. If there are no currently active line items, the most recent set of line items is returned.
	LineItems []LineItemInput `json:"line_items"`
	// Any additional attributes that you would like to save on the Order. The set of available options must be configured in advance in Lumos.
	CustomAttributes map[string]OrderCustomAttribute `json:"custom_attributes,omitempty"`
	// UUID of the application in Lumos where this order was sourced (e.g. the ID of Ironclad within Lumos)
	SourceAppID *string `json:"source_app_id,omitempty"`
}

func (o OrderInput) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSON(o, "", false)
}

func (o *OrderInput) UnmarshalJSON(data []byte) error {
	if err := utils.UnmarshalJSON(data, &o, "", false, false); err != nil {
		return err
	}
	return nil
}

func (o *OrderInput) GetUniqueIdentifier() string {
	if o == nil {
		return ""
	}
	return o.UniqueIdentifier
}

func (o *OrderInput) GetVendor() OrderInputVendor {
	if o == nil {
		return OrderInputVendor{}
	}
	return o.Vendor
}

func (o *OrderInput) GetStartDate() types.Date {
	if o == nil {
		return types.Date{}
	}
	return o.StartDate
}

func (o *OrderInput) GetEndDate() types.Date {
	if o == nil {
		return types.Date{}
	}
	return o.EndDate
}

func (o *OrderInput) GetOptOutDate() *types.Date {
	if o == nil {
		return nil
	}
	return o.OptOutDate
}

func (o *OrderInput) GetAutoRenewal() bool {
	if o == nil {
		return false
	}
	return o.AutoRenewal
}

func (o *OrderInput) GetCurrency() *string {
	if o == nil {
		return nil
	}
	return o.Currency
}

func (o *OrderInput) GetLineItems() []LineItemInput {
	if o == nil {
		return []LineItemInput{}
	}
	return o.LineItems
}

func (o *OrderInput) GetCustomAttributes() map[string]OrderCustomAttribute {
	if o == nil {
		return nil
	}
	return o.CustomAttributes
}

func (o *OrderInput) GetSourceAppID() *string {
	if o == nil {
		return nil
	}
	return o.SourceAppID
}
