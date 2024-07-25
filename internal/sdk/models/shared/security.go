// Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.

package shared

type Security struct {
	HTTPBearer *string `security:"scheme,type=http,subtype=bearer,name=Authorization"`
}

func (o *Security) GetHTTPBearer() *string {
	if o == nil {
		return nil
	}
	return o.HTTPBearer
}
