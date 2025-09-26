package table

import (
	"os"
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
	table.SetBorder(true)
	// This test mainly ensures the method doesn't panic
}

func TestSetCenterSeparator(t *testing.T) {
	table := New()
	table.SetCenterSeparator("|")
	// This test mainly ensures the method doesn't panic
}

func TestSetColumnSeparator(t *testing.T) {
	table := New()
	table.SetColumnSeparator("|")
	// This test mainly ensures the method doesn't panic
}

func TestSetRowSeparator(t *testing.T) {
	table := New()
	table.SetRowSeparator("-")
	// This test mainly ensures the method doesn't panic
}

func TestSetHeaderLine(t *testing.T) {
	table := New()
	table.SetHeaderLine(true)
	// This test mainly ensures the method doesn't panic
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
