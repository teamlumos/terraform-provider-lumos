package hooks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// appDeleteRecordKept handles DELETE /v1/apps/{id} answering 409 after it
// already disconnected an integration: the API keeps the record when other
// Lumos data references it. The API answers 409 for a refused delete too
// (an app not managed through the API), and the two differ only in their
// message, so the hook treats the delete as done only when the message says
// the record was kept and the app reads back disconnected. Any other 409
// stays an error, including a custom app's: the API's reset doesn't change a
// custom app, so when its record is kept nothing happened at all, App Store
// listing included.
//
// It is an AfterSuccess hook because the client runs those for every
// response the server sends, whatever its status code.
type appDeleteRecordKept struct{}

const appDeleteRecordKeptMessage = "its record could not be deleted"

func (appDeleteRecordKept) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if hookCtx.OperationID != "v1DeleteApp" || res.StatusCode != http.StatusConflict || res.Request == nil {
		return res, nil
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil || !bytes.Contains(body, []byte(appDeleteRecordKeptMessage)) {
		return res, nil
	}

	get, err := http.NewRequestWithContext(hookCtx.Context, http.MethodGet, res.Request.URL.String(), nil)
	if err != nil {
		return res, nil
	}
	get.Header.Set("Authorization", res.Request.Header.Get("Authorization"))
	got, err := hookCtx.SDKConfiguration.Client.Do(get)
	if err != nil {
		return res, nil
	}
	defer got.Body.Close()

	var app struct {
		Disconnected bool `json:"disconnected"`
	}
	if got.StatusCode != http.StatusOK || json.NewDecoder(got.Body).Decode(&app) != nil || !app.Disconnected {
		return res, nil
	}
	tflog.Warn(hookCtx.Context, "app was disconnected but kept because other Lumos data references it")
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Status:     "204 No Content",
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    res.Request,
	}, nil
}
