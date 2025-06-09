package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type tmpFilesRes struct {
	Data struct {
		URL string `json:"url"`
	} `json:"data"`
}

const (
	tmpFilesAPIURL          = "https://tmpfiles.org/api/v1/upload"
	tmpFilesResURL          = "tmpfiles.org/"
	tmpFilesURLWithDLSuffix = "tmpfiles.org/dl/"
	httpStr                 = "http://"
	httpsStr                = "https://"
)

func uploadImageToTmpFiles(imageBytes []byte, filename string) (string, error) {
	// Create a pipe to stream the multipart form data
	bodyReader, bodyWriter := io.Pipe()

	// Create a multipart writer
	writer := multipart.NewWriter(bodyWriter)

	// Run form creation in a goroutine to write to the pipe
	go func() {
		defer bodyWriter.Close()
		defer writer.Close()

		// Create a form file field named "file"
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			bodyWriter.CloseWithError(fmt.Errorf("failed to create form file: %w", err))
			return
		}

		// Write the image bytes to the form file field
		if _, err := io.Copy(part, bytes.NewReader(imageBytes)); err != nil {
			bodyWriter.CloseWithError(fmt.Errorf("failed to write image bytes: %w", err))
			return
		}
	}()

	req, err := http.NewRequest(http.MethodPost, tmpFilesAPIURL, bodyReader)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	fmt.Printf("[%d] %s\n", res.StatusCode, tmpFilesAPIURL)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var result tmpFilesRes
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	downloadURL := result.Data.URL
	unfixedURL := strings.Replace(downloadURL, tmpFilesResURL, tmpFilesURLWithDLSuffix, 1)
	httpsURL := strings.Replace(unfixedURL, httpStr, httpsStr, 1)

	return httpsURL, nil
}
