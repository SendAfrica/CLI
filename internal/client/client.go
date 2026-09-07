// Package client provides an HTTP client for the SendAfrica API.
// It wraps requests with the response envelope and auth headers.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

// Envelope is the standard SendAfrica response shape.
type Envelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   *APIError       `json:"error"`
}

// APIError represents an error from the API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

// Client is the SendAfrica API HTTP client.
type Client struct {
	BaseURL    string
	APIKey     string
	JWTToken   string
	HTTPClient *http.Client
}

// New creates a new API client.
func New(baseURL, apiKey, jwtToken string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		APIKey:     apiKey,
		JWTToken:   jwtToken,
		HTTPClient: &http.Client{},
	}
}

// RequestOpts configures a single request.
type RequestOpts struct {
	Method      string
	Path        string
	QueryParams map[string]string
	Body        interface{}
	Headers     map[string]string
}

// Do executes a request and returns the unwrapped envelope data.
func (c *Client) Do(opts RequestOpts) (json.RawMessage, error) {
	bodyReader, contentType, err := encodeBody(opts.Body)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL + opts.Path

	var req *http.Request
	if bodyReader != nil {
		req, err = http.NewRequest(opts.Method, url, bodyReader)
	} else {
		req, err = http.NewRequest(opts.Method, url, nil)
	}
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	// Auth: API key takes priority over JWT (a key is more specific to the CLI).
	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	} else if c.JWTToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.JWTToken)
	}

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	for k, v := range opts.QueryParams {
		q := req.URL.Query()
		q.Set(k, v)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

// DoWithResponse executes a request and returns the raw HTTP response.
func (c *Client) DoWithResponse(opts RequestOpts) (*http.Response, *Envelope, error) {
	bodyReader, contentType, err := encodeBody(opts.Body)
	if err != nil {
		return nil, nil, err
	}

	url := c.BaseURL + opts.Path

	var req *http.Request
	if bodyReader != nil {
		req, err = http.NewRequest(opts.Method, url, bodyReader)
	} else {
		req, err = http.NewRequest(opts.Method, url, nil)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("building request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", "application/json")

	if c.APIKey != "" {
		req.Header.Set("X-API-Key", c.APIKey)
	} else if c.JWTToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.JWTToken)
	}

	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}

	for k, v := range opts.QueryParams {
		q := req.URL.Query()
		q.Set(k, v)
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	env, err := parseEnvelope(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}

	return resp, env, nil
}

func (c *Client) parseResponse(resp *http.Response) (json.RawMessage, error) {
	env, err := parseEnvelope(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	if env.Error != nil {
		return nil, env.Error
	}
	if !env.Success {
		return nil, fmt.Errorf("request failed: success=false")
	}
	return env.Data, nil
}

func parseEnvelope(body io.Reader) (*Envelope, error) {
	var env Envelope
	dec := json.NewDecoder(body)
	if err := dec.Decode(&env); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &env, nil
}

func encodeBody(body interface{}) (*bytes.Reader, string, error) {
	if body == nil {
		return nil, "", nil
	}

	switch b := body.(type) {
	case *bytes.Buffer:
		return bytes.NewReader(b.Bytes()), "application/octet-stream", nil
	case ReaderWithContentType:
		return bytes.NewReader(b.Bytes()), b.ContentType(), nil
	default:
		data, err := json.Marshal(b)
		if err != nil {
			return nil, "", fmt.Errorf("marshalling body: %w", err)
		}
		return bytes.NewReader(data), "application/json", nil
	}
}

// ReaderWithContentType allows callers to provide a pre-encoded body with
// an explicit content type (used for multipart uploads).
type ReaderWithContentType interface {
	Bytes() []byte
	ContentType() string
}

// MultipartForm builds a multipart/form-data body for file uploads.
func MultipartForm(fields map[string]string, fileField, filename string, fileContent []byte, contentType string) (*bytes.Buffer, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for name, value := range fields {
		if err := writer.WriteField(name, value); err != nil {
			return nil, "", fmt.Errorf("writing field %s: %w", name, err)
		}
	}

	if fileContent != nil {
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fileField, filename))
		h.Set("Content-Type", contentType)
		part, err := writer.CreatePart(h)
		if err != nil {
			return nil, "", fmt.Errorf("creating file part: %w", err)
		}
		if _, err := part.Write(fileContent); err != nil {
			return nil, "", fmt.Errorf("writing file content: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("closing multipart writer: %w", err)
	}

	return body, writer.FormDataContentType(), nil
}
