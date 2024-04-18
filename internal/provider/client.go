package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type LumosAPIClient struct {
	HostURL    string
	HTTPClient *http.Client
	ApiToken   string
}

func NewLumosAPIClient(hostURL string, apiToken string) (*LumosAPIClient, error) {
	// TO-DO: error?
	return &LumosAPIClient{
		HTTPClient: &http.Client{},
		HostURL:    hostURL,
		ApiToken:   apiToken,
	}, nil

}

func (c *LumosAPIClient) MakeRequest(method string, endpoint string, requestBody interface{}, responseType lumosAPIResponse) (interface{}, error) {
	url := fmt.Sprintf("%s/%s", c.HostURL, endpoint)

	// Convert request body to JSON if any
	var body []byte
	if requestBody != nil {
		var err error
		body, err = json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
	}

	// Create a request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiToken))

	// Create an HTTP client
	client := &http.Client{}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("HTTP method %s failed with status %s and body %s", method, resp.Status, responseBody)
	}

	// Create a new instance of the response type using reflection.
	responseTypeInstance := reflect.New(reflect.TypeOf(responseType).Elem()).Interface()

	if err := json.Unmarshal(responseBody, responseTypeInstance); err != nil {
		return nil, fmt.Errorf("error parsing response %s -- Details: %s", string(responseBody), err)
	}

	return responseTypeInstance, nil

}

func (c *LumosAPIClient) createPermission(p requestablePermissionResourceModel) (*lumosAPIAppStoreRequestablePermissionResponse, error) {
	endpoint := PERMISSIONS_URL
	payload := buildAppStoreAppRequestablePermissionPayload(p)

	var reqPermission lumosAPIAppStoreRequestablePermissionResponse
	postResp, err := c.MakeRequest("POST", endpoint, payload, &reqPermission)
	if err != nil {
		return nil, err
	}
	permission, ok := postResp.(*lumosAPIAppStoreRequestablePermissionResponse)
	if !ok {
		return nil, err
	}
	return permission, nil
}

func (c *LumosAPIClient) getPermission(id string) (*lumosAPIAppStoreRequestablePermissionResponse, error) {
	endpoint := fmt.Sprintf(PERMISSION_BY_ID_URL, id)
	var reqPermission lumosAPIAppStoreRequestablePermissionResponse
	respInterface, err := c.MakeRequest("GET", endpoint, nil, &reqPermission)
	if err != nil {
		fmt.Printf("Error getting permission %s: %s", id, err)
		return nil, err
	}

	result, ok := respInterface.(*lumosAPIAppStoreRequestablePermissionResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	return result, nil
}

func (c *LumosAPIClient) updatePermission(id string, p requestablePermissionResourceModel) (*lumosAPIAppStoreRequestablePermissionResponse, error) {
	endpoint := fmt.Sprintf(PERMISSION_BY_ID_URL, id)
	payload := buildAppStoreAppRequestablePermissionPayload(p)

	var reqPermission lumosAPIAppStoreRequestablePermissionResponse
	postResp, err := c.MakeRequest("PATCH", endpoint, payload, &reqPermission)
	if err != nil {
		return nil, err
	}
	permission, ok := postResp.(*lumosAPIAppStoreRequestablePermissionResponse)
	if !ok {
		return nil, err
	}
	return permission, nil
}

func (c *LumosAPIClient) deletePermission(id string) error {
	endpoint := fmt.Sprintf(PERMISSION_BY_ID_URL, id)
	var reqPermission lumosAPIAppStoreRequestablePermissionResponse
	_, err := c.MakeRequest("DELETE", endpoint, nil, &reqPermission)
	if err != nil {
		fmt.Printf("Error deleting permission %s: %s", id, err)
		return err
	}

	return nil
}
