package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type CdnService struct{}

func (CdnService) Send(bucket, path, filename string, body io.Reader) (string, error) {
	var url string

	// CDN Token and Endpoint
	cdnUrl := os.Getenv("CDN_URL")
	cdnToken := os.Getenv("CDN_TOKEN")
	if cdnToken == "" || cdnUrl == "" {
		return "", errors.New("missing CDN configuration")
	}

	// Create multipart form
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	// Add file field
	filePart, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(filePart, body)
	if err != nil {
		return "", err
	}

	// Add other fields
	if err := writer.WriteField("bucket", bucket); err != nil {
		return "", err
	}
	if err := writer.WriteField("path", path); err != nil {
		return "", err
	}

	// Close the writer to finalize the multipart form
	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", cdnUrl+"upload", payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+cdnToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse JSON response
	var jsonResponse struct {
		Error bool   `json:"error"`
		Link  string `json:"link"`
	}
	if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
		return "", err
	}

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK || jsonResponse.Error {
		return "", errors.New("failed to upload file")
	}

	// Return the URL
	url = jsonResponse.Link
	return url, nil
}
