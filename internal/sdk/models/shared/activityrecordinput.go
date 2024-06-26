// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type ActivityRecordInput struct {
	// The activity records to upload.
	Records []ActivityRecord `json:"records"`
}

func (o *ActivityRecordInput) GetRecords() []ActivityRecord {
	if o == nil {
		return []ActivityRecord{}
	}
	return o.Records
}
