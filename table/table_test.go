package table

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/olekukonko/tablewriter/tw"
)

func TestNew(t *testing.T) {
	table := New()
	if table == nil {
		t.Error("New() returned nil")
	}
}

func TestNewWithWriter(t *testing.T) {
	table := NewWithWriter(os.Stdout)
	if table == nil {
		t.Error("NewWithWriter() returned nil")
	}
}

func TestSetHeaders(t *testing.T) {
	table := New()
	headers := []string{"Name", "Age", "City"}
	table.SetHeaders(headers)
	// This test mainly ensures the method doesn't panic
}

func TestSetHeadersConvertsToUppercase(t *testing.T) {
	table := New()
	headers := []string{"name", "age", "city"}
	table.SetHeaders(headers)

	// Check that headers are stored in uppercase internally
	if table.headers[0] != "NAME" {
		t.Errorf("Expected header 'NAME', got '%s'", table.headers[0])
	}
	if table.headers[1] != "AGE" {
		t.Errorf("Expected header 'AGE', got '%s'", table.headers[1])
	}
	if table.headers[2] != "CITY" {
		t.Errorf("Expected header 'CITY', got '%s'", table.headers[2])
	}
}

func TestAddRow(t *testing.T) {
	table := New()
	headers := []string{"Name", "Age"}
	table.SetHeaders(headers)

	row := []string{"John", "30"}
	table.AddRow(row)
	// This test mainly ensures the method doesn't panic
}

func TestAddRows(t *testing.T) {
	table := New()
	headers := []string{"Name", "Age"}
	table.SetHeaders(headers)

	rows := [][]string{
		{"John", "30"},
		{"Jane", "25"},
	}
	table.AddRows(rows)
	// This test mainly ensures the method doesn't panic
}

func TestSetBorder(t *testing.T) {
	table := New()

	// Test default is borderless
	if table.border {
		t.Error("Expected default border to be false")
	}

	// Test setting border to true
	table.SetBorder(true)
	if !table.border {
		t.Error("Expected border to be true after SetBorder(true)")
	}

	// Test setting border to false
	table.SetBorder(false)
	if table.border {
		t.Error("Expected border to be false after SetBorder(false)")
	}
}

func TestSetCenterSeparator(t *testing.T) {
	table := New()

	// Test default center separator
	if table.centerSeparator != "+" {
		t.Errorf("Expected default center separator to be '+', got '%s'", table.centerSeparator)
	}

	// Test setting center separator
	table.SetCenterSeparator("|")
	if table.centerSeparator != "|" {
		t.Errorf("Expected center separator to be '|', got '%s'", table.centerSeparator)
	}
}

func TestSetColumnSeparator(t *testing.T) {
	table := New()

	// Test default column separator
	if table.columnSeparator != "|" {
		t.Errorf("Expected default column separator to be '|', got '%s'", table.columnSeparator)
	}

	// Test setting column separator
	table.SetColumnSeparator("||")
	if table.columnSeparator != "||" {
		t.Errorf("Expected column separator to be '||', got '%s'", table.columnSeparator)
	}
}

func TestSetRowSeparator(t *testing.T) {
	table := New()

	// Test default row separator
	if table.rowSeparator != "-" {
		t.Errorf("Expected default row separator to be '-', got '%s'", table.rowSeparator)
	}

	// Test setting row separator
	table.SetRowSeparator("=")
	if table.rowSeparator != "=" {
		t.Errorf("Expected row separator to be '=', got '%s'", table.rowSeparator)
	}
}

func TestSetHeaderLine(t *testing.T) {
	table := New()

	// Test default header line
	if table.headerLine {
		t.Error("Expected default header line to be false")
	}

	// Test setting header line to true
	table.SetHeaderLine(true)
	if !table.headerLine {
		t.Error("Expected header line to be true after SetHeaderLine(true)")
	}

	// Test setting header line to false
	table.SetHeaderLine(false)
	if table.headerLine {
		t.Error("Expected header line to be false after SetHeaderLine(false)")
	}
}

func TestSetColumnAlignment(t *testing.T) {
	table := New()
	alignment := []tw.Align{tw.AlignLeft, tw.AlignCenter}
	table.SetColumnAlignment(alignment)
	// This test mainly ensures the method doesn't panic
}

func TestSetHeaderAlignment(t *testing.T) {
	table := New()
	alignment := []tw.Align{tw.AlignLeft, tw.AlignCenter}
	table.SetHeaderAlignment(alignment)
	// This test mainly ensures the method doesn't panic
}

func TestSetAutoWrapText(t *testing.T) {
	table := New()
	table.SetAutoWrapText(true)
	// This test mainly ensures the method doesn't panic
}

func TestSetReflowDuringAutoWrap(t *testing.T) {
	table := New()
	table.SetReflowDuringAutoWrap(true)
	// This test mainly ensures the method doesn't panic
}

func TestSetTablePadding(t *testing.T) {
	table := New()
	table.SetTablePadding("  ")
	// This test mainly ensures the method doesn't panic
}

func TestSetNoWhiteSpace(t *testing.T) {
	table := New()
	table.SetNoWhiteSpace(true)
	// This test mainly ensures the method doesn't panic
}

func TestSetFooter(t *testing.T) {
	table := New()
	footer := []string{"Total", "2"}
	table.SetFooter(footer)
	// This test mainly ensures the method doesn't panic
}

func TestBorderlessRenderingDefault(t *testing.T) {
	var buf bytes.Buffer
	table := NewWithWriter(&buf)

	table.SetHeaders([]string{"name", "age"})
	table.AddRow([]string{"John", "30"})
	table.AddRow([]string{"Jane", "25"})

	err := table.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}

	output := buf.String()
	// Should not contain border characters
	if strings.Contains(output, "|") || strings.Contains(output, "+") || strings.Contains(output, "-") {
		t.Errorf("Borderless table should not contain border characters, got:\n%s", output)
	}

	// Should contain headers in uppercase
	if !strings.Contains(output, "NAME") || !strings.Contains(output, "AGE") {
		t.Errorf("Expected uppercase headers NAME and AGE, got:\n%s", output)
	}
}

func TestBorderedRendering(t *testing.T) {
	var buf bytes.Buffer
	table := NewWithWriter(&buf)

	table.SetHeaders([]string{"name", "age"})
	table.AddRow([]string{"John", "30"})
	table.AddRow([]string{"Jane", "25"})
	table.SetBorder(true)

	err := table.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}

	output := buf.String()
	// Should contain border characters
	if !strings.Contains(output, "|") || !strings.Contains(output, "+") || !strings.Contains(output, "-") {
		t.Errorf("Bordered table should contain border characters, got:\n%s", output)
	}

	// Should contain headers in uppercase
	if !strings.Contains(output, "NAME") || !strings.Contains(output, "AGE") {
		t.Errorf("Expected uppercase headers NAME and AGE, got:\n%s", output)
	}
}

func TestBorderedRenderingWithHeaderLine(t *testing.T) {
	var buf bytes.Buffer
	table := NewWithWriter(&buf)

	table.SetHeaders([]string{"name", "age"})
	table.AddRow([]string{"John", "30"})
	table.SetBorder(true)
	table.SetHeaderLine(true)

	err := table.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Should have at least 4 lines: top border, header, header separator, data row, bottom border
	if len(lines) < 5 {
		t.Errorf("Expected at least 5 lines with header line enabled, got %d:\n%s", len(lines), output)
	}
}

func TestCustomBorderCharacters(t *testing.T) {
	var buf bytes.Buffer
	table := NewWithWriter(&buf)

	table.SetHeaders([]string{"name", "age"})
	table.AddRow([]string{"John", "30"})
	table.SetBorder(true)
	table.SetCenterSeparator("*")
	table.SetColumnSeparator(":")
	table.SetRowSeparator("=")

	err := table.Render()
	if err != nil {
		t.Errorf("Render failed: %v", err)
	}

	output := buf.String()
	// Should contain custom border characters
	if !strings.Contains(output, "*") || !strings.Contains(output, ":") || !strings.Contains(output, "=") {
		t.Errorf("Expected custom border characters *, :, =, got:\n%s", output)
	}
}
