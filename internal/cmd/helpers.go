package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cameltech/sendafrica-cli/internal/api"
)

// decodeData unmarshals json.RawMessage into the target.
func decodeData(data json.RawMessage, target interface{}) error {
	if len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// decodePaginated unmarshals a paginated API response, extracting the items
// array into target. The raw data is expected to be an object like
// {"items": [...], "page": 1, "total": 42}.
func decodePaginated(data json.RawMessage, target interface{}) (api.PaginatedResponse, error) {
	var page api.PaginatedResponse
	if err := json.Unmarshal(data, &page); err != nil {
		return page, fmt.Errorf("decoding paginated response: %w", err)
	}
	if err := page.ExtractItems(target); err != nil {
		return page, fmt.Errorf("extracting items: %w", err)
	}
	return page, nil
}

// decodeListOrPaginated tries a bare-array decode first, then falls back to
// a paginated {items: [...]} decode. This handles APIs where some list
// endpoints return arrays and others return paginated objects.
func decodeListOrPaginated(data json.RawMessage, target interface{}) error {
	// Try bare array first.
	if err := json.Unmarshal(data, target); err == nil {
		return nil
	}
	// Fall back to paginated object.
	_, err := decodePaginated(data, target)
	return err
}

// rawMessage converts json.RawMessage to a generic interface{} for printing.
func rawMessage(data json.RawMessage) interface{} {
	if len(data) == 0 {
		return nil
	}
	var v interface{}
	_ = json.Unmarshal(data, &v)
	return v
}

func jsonUnmarshal(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}

// parseStringFlag splits a comma-separated flag value into a slice.
func parseStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

// validateAuth returns an error if no API key or JWT token is configured.
func requireAuth() error {
	return ensureAuth()
}

// paginatedParams builds query params for paginated endpoints.
func paginatedParams(page, perPage int) map[string]string {
	params := make(map[string]string)
	if page > 0 {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if perPage > 0 {
		params["per_page"] = fmt.Sprintf("%d", perPage)
	}
	return params
}
