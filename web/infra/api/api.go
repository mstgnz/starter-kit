package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/mstgnz/starter-kit/web/infra/config"
	"github.com/mstgnz/starter-kit/web/model"
)

type ApiService struct {
	url         string
	path        string
	method      string
	params      map[string]any
	headers     map[string]string
	attachments []Attachment
}

type Attachment struct {
	Name     string
	Contents []byte
	Filename string
}

// New creates a new ApiService instance
func New() *ApiService {
	return &ApiService{
		url:         os.Getenv("API_URL"),
		params:      make(map[string]any),
		headers:     make(map[string]string),
		attachments: make([]Attachment, 0),
	}
}

func (r *ApiService) WithHeader(headers map[string]string) *ApiService {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

func (r *ApiService) WithToken(token string) *ApiService {
	r.headers["Authorization"] = "Bearer " + token
	return r
}

func (r *ApiService) WithAttachment(name string, contents []byte, filename string) *ApiService {
	r.attachments = append(r.attachments, Attachment{
		Name:     name,
		Contents: contents,
		Filename: filename,
	})
	return r
}

func (r *ApiService) Get(path string, query map[string]any) (*model.Response, error) {
	return r.prepare(http.MethodGet, path, query)
}

func (r *ApiService) Delete(path string, query map[string]any) (*model.Response, error) {
	return r.prepare(http.MethodDelete, path, query)
}

func (r *ApiService) Post(path string, data map[string]any) (*model.Response, error) {
	return r.prepare(http.MethodPost, path, data)
}

func (r *ApiService) Put(path string, data map[string]any) (*model.Response, error) {
	return r.prepare(http.MethodPut, path, data)
}

func (r *ApiService) prepare(method, path string, data map[string]any) (*model.Response, error) {
	r.path = path
	r.method = method

	// Clear existing parameters
	r.params = make(map[string]any)

	for k, v := range data {
		r.params[k] = v
	}
	r.params["lang"] = config.App().Lang

	return r.send()
}

func (r *ApiService) setHeaderHash(path string) error {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	secret := os.Getenv("APP_SECRET")
	rawData := fmt.Sprintf("Starter.%s:%s:%s.Kit", timestamp, path, secret)
	hash := sha256.Sum256([]byte(rawData))

	r.headers["Timestamp"] = timestamp
	r.headers["Hash"] = hex.EncodeToString(hash[:])

	return nil
}

func (r *ApiService) send() (*model.Response, error) {
	if err := r.setHeaderHash(r.path); err != nil {
		return nil, fmt.Errorf("hash generation error: %w", err)
	}

	if r.headers["Authorization"] == "" {
		r.headers["Authorization"] = "Bearer " + config.App().Token
	}

	url, reqBody, err := r.buildRequest()
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	req, err := http.NewRequest(r.method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("HTTP request generation error: %w", err)
	}

	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %w", err)
	}

	var response model.Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}

	// Insert CURL output
	if response.Data == nil {
		response.Data = make(map[string]any)
	}
	response.Data["curl"] = r.ToCurl()

	r.reset()
	return &response, nil
}

func (r *ApiService) buildRequest() (string, io.Reader, error) {
	baseURL := r.url + r.path

	switch r.method {
	case http.MethodGet, http.MethodDelete:
		params := url.Values{}
		for k, v := range r.params {
			params.Add(k, fmt.Sprintf("%v", v))
		}
		if len(params) > 0 {
			baseURL += "?" + params.Encode()
		}
		return baseURL, nil, nil

	case http.MethodPost, http.MethodPut:
		if len(r.attachments) > 0 {
			return r.buildMultipartRequest(baseURL)
		}
		return r.buildJSONRequest(baseURL)

	default:
		return "", nil, fmt.Errorf("unsupported HTTP method: %s", r.method)
	}
}

func (r *ApiService) buildMultipartRequest(baseURL string) (string, io.Reader, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, attachment := range r.attachments {
		part, err := writer.CreateFormFile(attachment.Name, attachment.Filename)
		if err != nil {
			return "", nil, fmt.Errorf("form file creation error: %w", err)
		}
		if _, err := part.Write(attachment.Contents); err != nil {
			return "", nil, fmt.Errorf("file write error: %w", err)
		}
	}

	for k, v := range r.params {
		if err := writer.WriteField(k, fmt.Sprintf("%v", v)); err != nil {
			return "", nil, fmt.Errorf("form field write error: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return "", nil, fmt.Errorf("multipart form closing error: %w", err)
	}

	r.headers["Content-Type"] = writer.FormDataContentType()
	return baseURL, body, nil
}

func (r *ApiService) buildJSONRequest(baseURL string) (string, io.Reader, error) {
	if len(r.params) == 0 {
		return baseURL, nil, nil
	}

	jsonData, err := json.Marshal(r.params)
	if err != nil {
		return "", nil, fmt.Errorf("JSON conversion error: %w", err)
	}

	r.headers["Content-Type"] = "application/json"
	return baseURL, bytes.NewBuffer(jsonData), nil
}

func (r *ApiService) reset() {
	r.path = ""
	r.method = ""
	r.params = make(map[string]any)
	r.headers = make(map[string]string)
	r.attachments = make([]Attachment, 0)
}

// ToCurl returns the curl command representation of the request
func (r *ApiService) ToCurl() string {
	var curl bytes.Buffer
	curl.WriteString("curl -X '" + r.method + "' \\\n")

	// Base URL
	baseURL := r.url + r.path

	// Add query parameters for GET requests
	if r.method == http.MethodGet && len(r.params) > 0 {
		baseURL += "?"
		first := true
		for k, v := range r.params {
			if !first {
				baseURL += "&"
			}
			baseURL += k + "=" + fmt.Sprintf("%v", v)
			first = false
		}
	}
	curl.WriteString("  '" + baseURL + "' \\\n")

	// Add headers
	for k, v := range r.headers {
		curl.WriteString("  -H '" + k + ": " + v + "' \\\n")
	}

	// Add body for POST/PUT requests
	if (r.method == http.MethodPost || r.method == http.MethodPut) && len(r.params) > 0 {
		jsonData, _ := json.Marshal(r.params)
		curl.WriteString("  -H 'Content-Type: application/json' \\\n")
		curl.WriteString("  -d '" + string(jsonData) + "'")
	}

	return curl.String()
}
