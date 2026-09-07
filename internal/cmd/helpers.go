package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
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
