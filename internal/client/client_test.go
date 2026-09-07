package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientWithAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "SA-testkey" {
			t.Errorf("expected X-API-Key header, got %q", r.Header.Get("X-API-Key"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"status": "ok"},
			"error":   nil,
		})
	}))
	defer server.Close()

	c := New(server.URL, "SA-testkey", "")
	data, err := c.Do(RequestOpts{Method: "GET", Path: "/v1/packages"})
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected status ok, got %s", result["status"])
	}
}

func TestClientWithJWT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-jwt" {
			t.Errorf("expected Authorization Bearer header, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-API-Key") != "" {
			t.Error("should not set X-API-Key when using JWT")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"name": "test"},
			"error":   nil,
		})
	}))
	defer server.Close()

	c := New(server.URL, "", "test-jwt")
	data, err := c.Do(RequestOpts{Method: "GET", Path: "/v1/auth/me"})
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if result["name"] != "test" {
		t.Errorf("expected name test, got %s", result["name"])
	}
}

func TestAPIKeyTakesPrecedenceOverJWT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "SA-api-key" {
			t.Errorf("expected API key to take precedence, got X-API-Key=%q", r.Header.Get("X-API-Key"))
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("should not set Authorization when API key is present")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    nil,
			"error":   nil,
		})
	}))
	defer server.Close()

	c := New(server.URL, "SA-api-key", "jwt-token")
	_, err := c.Do(RequestOpts{Method: "GET", Path: "/test"})
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
}

func TestAPIErrorReturned(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"data":    nil,
			"error": map[string]string{
				"code":    "authentication_failed",
				"message": "Invalid API key",
			},
		})
	}))
	defer server.Close()

	c := New(server.URL, "SA-badkey", "")
	_, err := c.Do(RequestOpts{Method: "GET", Path: "/test"})
	if err == nil {
		t.Fatal("expected error for failed auth")
	}
	if !strings.Contains(err.Error(), "authentication_failed") {
		t.Errorf("expected authentication_failed in error, got %q", err.Error())
	}
}

func TestQueryParams(t *testing.T) {
	var capturedQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    []string{},
			"error":   nil,
		})
	}))
	defer server.Close()

	c := New(server.URL, "key", "")
	_, err := c.Do(RequestOpts{
		Method: "GET",
		Path:   "/v1/sms/logs",
		QueryParams: map[string]string{
			"page":     "2",
			"per_page": "100",
		},
	})
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}

	if !strings.Contains(capturedQuery, "page=2") {
		t.Errorf("expected page=2 in query, got %q", capturedQuery)
	}
	if !strings.Contains(capturedQuery, "per_page=100") {
		t.Errorf("expected per_page=100 in query, got %q", capturedQuery)
	}
}

func TestPOSTWithBody(t *testing.T) {
	var receivedBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&receivedBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    map[string]string{"id": "123"},
			"error":   nil,
		})
	}))
	defer server.Close()

	c := New(server.URL, "key", "")
	reqBody := map[string]string{"to": "+255712345678", "message": "test"}
	data, err := c.Do(RequestOpts{
		Method: "POST",
		Path:   "/v1/sms/",
		Body:   reqBody,
	})
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}

	if receivedBody["to"] != "+255712345678" {
		t.Errorf("expected to in body, got %v", receivedBody["to"])
	}
	if receivedBody["message"] != "test" {
		t.Errorf("expected message in body, got %v", receivedBody["message"])
	}

	var result map[string]string
	json.Unmarshal(data, &result)
	if result["id"] != "123" {
		t.Errorf("expected id 123, got %s", result["id"])
	}
}

func TestMultipartForm(t *testing.T) {
	body, contentType, err := MultipartForm(
		map[string]string{"field1": "value1"},
		"file",
		"test.csv",
		[]byte("col1,col2\nval1,val2\n"),
		"text/csv",
	)
	if err != nil {
		t.Fatalf("MultipartForm failed: %v", err)
	}
	if !strings.Contains(contentType, "multipart/form-data") {
		t.Errorf("expected multipart/form-data, got %s", contentType)
	}
	bodyStr := body.String()
	if !strings.Contains(bodyStr, "field1") {
		t.Error("expected field1 in body")
	}
	if !strings.Contains(bodyStr, "test.csv") {
		t.Error("expected filename in body")
	}
	if !strings.Contains(bodyStr, "col1,col2") {
		t.Error("expected file content in body")
	}
}
