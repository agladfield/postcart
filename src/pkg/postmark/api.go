// Package postmark wraps calling the Postmark API and provides request and response types
package postmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/agladfield/postcart/pkg/shared/env"
)

const (
	apiURL               = "https://api.postmarkapp.com"
	serverTokenHeaderKey = "X-Postmark-Server-Token"
	acceptsHeaderKey     = "Accepts"
	contentTypeKey       = "Content-Type"
	applicationJSONValue = "application/json"
)

const (
	postmarkAPIErrFmtStr        = "postmark err: %w"
	postmarkAPIConfigErrFmtStr  = "postmark configuration err: %w"
	postmarkAPIRequestErrFmtStr = "postmark request err: %w"
)

// POSTMARK_API_TEST can use this as server token for testing validity but not send

type apiClient struct {
	baseURL        string
	client         *http.Client
	defaultHeaders *http.Header
}

var api *apiClient

func Configure() error {
	// validate keys or whatever
	// validation errors happen here

	postmarkClient := &http.Client{
		Timeout: time.Second * 5,
	}

	defaultHeaders := make(http.Header)
	defaultHeaders.Add(serverTokenHeaderKey, env.PostmarkServerToken())
	defaultHeaders.Add(acceptsHeaderKey, applicationJSONValue)

	api = &apiClient{
		baseURL:        apiURL,
		client:         postmarkClient,
		defaultHeaders: &defaultHeaders,
	}

	return nil
}

// DecodeToStruct simplifies decoding the raw postmark requests and responses into their structured data
func DecodeToStruct[T any](body io.ReadCloser, ptr *T) error {
	decodeErr := json.NewDecoder(body).Decode(ptr)
	if decodeErr != nil {
		return fmt.Errorf("failed to decode body to struct: %w", decodeErr)
	}

	return nil
}

// EncodeToStruct simplifies encoding the structured postmark requests into a raw format suitable for http requests
func EncodeToStruct(data any) (*bytes.Buffer, error) {
	jsonBytes, marshallErr := json.Marshal(data)
	if marshallErr != nil {
		return nil, fmt.Errorf("failed to encode struct to body: %w", marshallErr)
	}
	byteBuffer := bytes.NewBuffer(jsonBytes)

	return byteBuffer, nil
}

func request[R any](c *apiClient, path string, method string, body io.Reader, resData *R) error {
	joinedURL := fmt.Sprintf("%s%s", c.baseURL, path)

	reqHeaders := c.defaultHeaders.Clone()
	var req *http.Request
	var newReqErr error

	if body != nil {
		req, newReqErr = http.NewRequest(method, joinedURL, body)
		reqHeaders.Add(contentTypeKey, applicationJSONValue)
	} else {
		req, newReqErr = http.NewRequest(method, joinedURL, nil)
	}
	if newReqErr != nil {
		return fmt.Errorf(postmarkAPIRequestErrFmtStr, fmt.Errorf("api client request creation err: %w", newReqErr))
	}
	req.Header = reqHeaders

	res, resErr := c.client.Do(req)
	fmt.Printf("[%d] %s\n", res.StatusCode, joinedURL)
	if resErr != nil {
		return fmt.Errorf(postmarkAPIRequestErrFmtStr, fmt.Errorf("api client response err: %w [%d]", resErr, res.StatusCode))
	}
	defer res.Body.Close()

	decodeErr := DecodeToStruct(res.Body, resData)
	if decodeErr != nil {
		return fmt.Errorf(postmarkAPIRequestErrFmtStr, (fmt.Errorf("api client failed to decode response err: %w", decodeErr)))
	}

	// return fmt.Errorf("res data: %v", *resData)

	return nil
}
