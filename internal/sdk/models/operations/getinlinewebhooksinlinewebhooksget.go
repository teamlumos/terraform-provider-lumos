// Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.

package operations

import (
	"github.com/teamlumos/terraform-provider-lumos/internal/sdk/models/shared"
	"net/http"
)

type GetInlineWebhooksInlineWebhooksGetResponse struct {
	// HTTP response content type for this operation
	ContentType string
	// HTTP response status code for this operation
	StatusCode int
	// Raw HTTP response; suitable for custom response parsing
	RawResponse *http.Response
	// Successful Response
	ResponseGetInlineWebhooksInlineWebhooksGet []shared.InlineWebhook
}

func (o *GetInlineWebhooksInlineWebhooksGetResponse) GetContentType() string {
	if o == nil {
		return ""
	}
	return o.ContentType
}

func (o *GetInlineWebhooksInlineWebhooksGetResponse) GetStatusCode() int {
	if o == nil {
		return 0
	}
	return o.StatusCode
}

func (o *GetInlineWebhooksInlineWebhooksGetResponse) GetRawResponse() *http.Response {
	if o == nil {
		return nil
	}
	return o.RawResponse
}

func (o *GetInlineWebhooksInlineWebhooksGetResponse) GetResponseGetInlineWebhooksInlineWebhooksGet() []shared.InlineWebhook {
	if o == nil {
		return nil
	}
	return o.ResponseGetInlineWebhooksInlineWebhooksGet
}