package cmd

import "testing"

func TestMask(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"short", "short", "****"},
		{"exact-8", "12345678", "****"},
		{"long", "SA-abcdefghijklmnopqrstuvwxyz1234567890", "SA-a****7890"},
		{"medium", "1234567890", "1234****7890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mask(tt.input)
			if got != tt.want {
				t.Errorf("mask(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseStringSlice(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", nil},
		{"a", []string{"a"}},
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b , c ", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseStringSlice(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("parseStringSlice(%q) = %v (len %d), want %v (len %d)",
					tt.input, got, len(got), tt.want, len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseStringSlice(%q)[%d] = %q, want %q",
						tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}
