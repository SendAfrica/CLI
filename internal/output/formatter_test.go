package output

import (
	"bytes"
	"strings"
	"testing"
)

type testItem struct {
	ID    string `json:"id" table:"ID"`
	Name  string `json:"name" table:"Name"`
	Value int    `json:"value" table:"Value"`
}

func TestPrintJSON(t *testing.T) {
	p := New(FormatJSON)
	data := testItem{ID: "1", Name: "test", Value: 42}
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id": "1"`) {
		t.Errorf("expected id in output, got: %s", out)
	}
	if !strings.Contains(out, `"name": "test"`) {
		t.Errorf("expected name in output, got: %s", out)
	}
}

func TestPrintJSONSlice(t *testing.T) {
	p := New(FormatJSON)
	data := []testItem{
		{ID: "1", Name: "Alice", Value: 10},
		{ID: "2", Name: "Bob", Value: 20},
	}
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alice") {
		t.Errorf("expected Alice in output, got: %s", out)
	}
	if !strings.Contains(out, "Bob") {
		t.Errorf("expected Bob in output, got: %s", out)
	}
}

func TestPrintTableSingleStruct(t *testing.T) {
	p := New(FormatTable)
	data := testItem{ID: "1", Name: "test", Value: 42}
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "1") {
		t.Errorf("expected ID and value in table output, got: %s", out)
	}
	if !strings.Contains(out, "Name") || !strings.Contains(out, "test") {
		t.Errorf("expected Name and test in table output, got: %s", out)
	}
	if !strings.Contains(out, "Value") || !strings.Contains(out, "42") {
		t.Errorf("expected Value and 42 in table output, got: %s", out)
	}
}

func TestPrintTableSlice(t *testing.T) {
	p := New(FormatTable)
	data := []testItem{
		{ID: "1", Name: "Alice", Value: 10},
		{ID: "2", Name: "Bob", Value: 20},
	}
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "Name") || !strings.Contains(out, "Value") {
		t.Errorf("expected column headers, got: %s", out)
	}
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("expected item names in output, got: %s", out)
	}
}

func TestPrintTableEmptySlice(t *testing.T) {
	p := New(FormatTable)
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print([]testItem{}); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "No results." {
		t.Errorf("expected 'No results.', got: %s", out)
	}
}

func TestPrintYAML(t *testing.T) {
	p := New(FormatYAML)
	data := testItem{ID: "1", Name: "test", Value: 42}
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "id: \"1\"") {
		t.Errorf("expected id in YAML, got: %s", out)
	}
	if !strings.Contains(out, "name: test") {
		t.Errorf("expected name in YAML, got: %s", out)
	}
}

func TestDefaultFormatIsTable(t *testing.T) {
	p := New("")
	if p.Format() != FormatTable {
		t.Errorf("expected default format table, got %s", p.Format())
	}
}

func TestNilData(t *testing.T) {
	p := New(FormatTable)
	var buf bytes.Buffer
	p2 := p.WithWriter(&buf)
	if err := p2.Print(nil); err != nil {
		t.Fatalf("Print of nil should not fail: %v", err)
	}
}
