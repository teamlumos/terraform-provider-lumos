package hooks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// approvalChainKept treats DELETE /v1/workflows/{id} answering 422
// workflow.approval_chain_delete_not_supported as done. Lumos never deletes
// an approval chain, even one no template uses, so without this a
// lumos_workflow of type approval-chain could never be destroyed. The chain
// stays in Lumos and leaves Terraform state; lumos_workflow's plan warns
// about that before the apply. Any other refusal stays an error.
//
// It is an AfterSuccess hook because the client runs those for every
// response the server sends, whatever its status code.
type approvalChainKept struct{}

const approvalChainDeleteNotSupported = "workflow.approval_chain_delete_not_supported"

func (approvalChainKept) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if hookCtx.OperationID != "deleteWorkflow" || res.StatusCode != http.StatusUnprocessableEntity {
		return res, nil
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	res.Body = io.NopCloser(bytes.NewReader(body))
	if err != nil {
		return res, nil
	}
	var refusal struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &refusal) != nil || refusal.Error.Code != approvalChainDeleteNotSupported {
		return res, nil
	}
	tflog.Warn(hookCtx.Context, "Lumos keeps approval chains; the chain was removed from state but stays in Lumos")
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Status:     "204 No Content",
		Header:     http.Header{},
		Body:       http.NoBody,
		Request:    res.Request,
	}, nil
}
