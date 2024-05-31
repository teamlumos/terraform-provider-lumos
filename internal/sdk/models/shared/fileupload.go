// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package shared

type FileUpload struct {
	// The type of content being uploaded. The only supported type today is 'url'
	Type FileUploadType `json:"type"`
	// The data to upload. This should be a publicly accessible URL that resolves to a file
	Value string `json:"value"`
}

func (o *FileUpload) GetType() FileUploadType {
	if o == nil {
		return FileUploadType("")
	}
	return o.Type
}

func (o *FileUpload) GetValue() string {
	if o == nil {
		return ""
	}
	return o.Value
}
