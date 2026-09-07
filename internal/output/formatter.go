// Package output provides formatters for CLI command output.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

// Format constants.
const (
	FormatTable = "table"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
)

// Printer wraps output format selection and writes to stdout.
type Printer struct {
	format string
	w      io.Writer
}

func New(format string) *Printer {
	if format == "" {
		format = FormatTable
	}
	return &Printer{format: format, w: os.Stdout}
}

// WithWriter returns a new Printer writing to w (useful for testing).
func (p *Printer) WithWriter(w io.Writer) *Printer {
	return &Printer{format: p.format, w: w}
}

// Format returns the configured format string.
func (p *Printer) Format() string {
	return p.format
}

// Print writes data in the configured format.
func (p *Printer) Print(data interface{}) error {
	return p.print(p.w, data)
}

func (p *Printer) print(w io.Writer, data interface{}) error {
	switch p.format {
	case FormatJSON:
		return p.printJSON(w, data)
	case FormatYAML:
		return p.printYAML(w, data)
	default:
		return p.printTable(w, data)
	}
}

func (p *Printer) printJSON(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	return nil
}

func (p *Printer) printYAML(w io.Writer, data interface{}) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("encoding YAML: %w", err)
	}
	return nil
}

func (p *Printer) printTable(w io.Writer, data interface{}) error {
	val := data
	if val == nil {
		return nil
	}

	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return printSliceTable(w, rv)
	default:
		printKeyValueTable(w, data)
		return nil
	}
}

// printSliceTable prints a slice of structs as a column-based table.
func printSliceTable(w io.Writer, rv reflect.Value) error {
	if rv.Len() == 0 {
		fmt.Fprintln(w, "No results.")
		return nil
	}

	first := rv.Index(0)
	for first.Kind() == reflect.Ptr {
		if first.IsNil() {
			return nil
		}
		first = first.Elem()
	}

	if first.Kind() != reflect.Struct {
		fmt.Fprintf(w, "%v\n", rv.Index(0).Interface())
		return nil
	}

	fields := structFields(first.Type())

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	for _, f := range fields {
		fmt.Fprintf(tw, "%s\t", f.name)
	}
	fmt.Fprintln(tw)
	for range fields {
		fmt.Fprint(tw, "---\t")
	}
	fmt.Fprintln(tw)

	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		for elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				elem = reflect.Value{}
				break
			}
			elem = elem.Elem()
		}
		for _, f := range fields {
			if elem.IsValid() {
				fv := elem.Field(f.index)
				fmt.Fprintf(tw, "%s\t", formatValue(fv))
			} else {
				fmt.Fprint(tw, "\t")
			}
		}
		fmt.Fprintln(tw)
	}

	return tw.Flush()
}

// printKeyValueTable prints a single struct as key/value pairs.
func printKeyValueTable(w io.Writer, data interface{}) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	rv := reflect.ValueOf(data)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		fmt.Fprintf(w, "%v\n", data)
		return
	}

	fields := structFields(rv.Type())
	for _, f := range fields {
		fv := rv.Field(f.index)
		fmt.Fprintf(tw, "%s\t%s\n", f.name, formatValue(fv))
	}
	tw.Flush()
}

type fieldInfo struct {
	name  string
	index int
}

func structFields(rt reflect.Type) []fieldInfo {
	var fields []fieldInfo
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if field.Anonymous {
			continue
		}
		tag := field.Tag.Get("table")
		if tag == "-" {
			continue
		}
		name := field.Name
		if tag != "" && tag != "-" {
			name = tag
		}
		fields = append(fields, fieldInfo{name: name, index: i})
	}
	return fields
}

func formatValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	if !v.CanInterface() {
		return ""
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return ""
		}
		return formatValue(v.Elem())
	case reflect.Struct:
		if b, err := json.Marshal(v.Interface()); err == nil {
			return string(b)
		}
		return fmt.Sprintf("%v", v.Interface())
	case reflect.Slice, reflect.Array:
		if v.Len() == 0 {
			return ""
		}
		if b, err := json.Marshal(v.Interface()); err == nil {
			return string(b)
		}
		return fmt.Sprintf("%v", v.Interface())
	case reflect.Map:
		if v.Len() == 0 {
			return ""
		}
		if b, err := json.Marshal(v.Interface()); err == nil {
			return string(b)
		}
		return fmt.Sprintf("%v", v.Interface())
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	case reflect.Interface:
		if v.IsNil() {
			return ""
		}
		return formatValue(reflect.ValueOf(v.Interface()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
