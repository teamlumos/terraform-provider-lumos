package hooks

import (
	"net/http"
	"strings"
)

type RequestablePermissionsHook struct{}

var (
	_ sdkInitHook       = (*RequestablePermissionsHook)(nil)
	_ beforeRequestHook = (*RequestablePermissionsHook)(nil)
	_ afterSuccessHook  = (*RequestablePermissionsHook)(nil)
	_ afterErrorHook    = (*RequestablePermissionsHook)(nil)
)

func (i *RequestablePermissionsHook) SDKInit(baseURL string, client HTTPClient) (string, HTTPClient) {
	// modify the baseURL or wrap the client used by the SDK here and return the updated values
	return baseURL, client
}

func (i *RequestablePermissionsHook) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	// modify the request object before it is sent, such as adding headers or query parameters or return an error to stop the request from being sent
	if strings.Contains(req.URL.String(), "requestable_permissions") {
		q := req.URL.Query()
		q.Add("include_inherited_configs", "false")
		req.URL.RawQuery = q.Encode()
	}

	return req, nil
}

func (i *RequestablePermissionsHook) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	// modify the response object before deserialization or return an error to stop the response from being deserialized
	return res, nil
}

func (i *RequestablePermissionsHook) AfterError(hookCtx AfterErrorContext, res *http.Response, err error) (*http.Response, error) {
	// modify the response before it is deserialized as a custom error or the error object before it is returned or return an error wrapped in the FailEarly error in this package to exit from the hook chain early
	return res, err
}
