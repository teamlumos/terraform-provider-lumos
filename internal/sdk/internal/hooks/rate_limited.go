package hooks

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// rateLimited turns a 429 that outlasted the generated retries into an error
// that says the request was rate limited. The generated retry loop (the
// overlay's root x-speakeasy-retries) waits out each Retry-After, but gives
// up at once when Retry-After runs past its maxElapsedTime, and then returns
// the 429 response without an error, which the provider reports only as an
// unexpected status code.
type rateLimited struct{}

func (rateLimited) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if res.StatusCode != http.StatusTooManyRequests {
		return res, nil
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
	retryAfter := res.Header.Get("Retry-After")
	if retryAfter == "" {
		retryAfter = "not sent"
	}
	// "2 minutes" is maxElapsedTime in overlay.yaml.
	return nil, fmt.Errorf(
		"%s: the Lumos API is still rate limiting this request (status 429, Retry-After: %s) after the provider retried it for up to 2 minutes. Re-run later, or lower terraform -parallelism. Response: %s",
		hookCtx.OperationID, retryAfter, bytes.TrimSpace(body),
	)
}
