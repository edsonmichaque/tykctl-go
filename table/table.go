package table

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/edsonmichaque/tykctl-go/terminal"
)

// Table represents a table for formatted output
type Table struct {
	headers    []string
	rows       [][]string
	output     io.Writer
	terminal   *terminal.Terminal
	separator  string
	alignment  []int
	widths     []int
}

// New creates a new table instance
func New() *Table {
	return &Table{
		output:    os.Stdout,
		terminal:  terminal.New(),
		separator: "\t",
	}
}

// NewWithWriter creates a new table with a custom writer
func NewWithWriter(w io.Writer) *Table {
	return &Table{
		output:    w,
		terminal:  terminal.New(),
		separator: "\t",
	}
}

// SetHeaders sets the table headers
func (t *Table) SetHeaders(headers []string) {
	t.headers = headers
	t.widths = make([]int, len(headers))
	for i, header := range headers {
		t.widths[i] = len(header)
	}
}

// AddRow adds a row to the table
func (t *Table) AddRow(row []string) {
	t.rows = append(t.rows, row)
	
	// Update column widths
	for i, cell := range row {
		if i < len(t.widths) {
			if len(cell) > t.widths[i] {
				t.widths[i] = len(cell)
			}
		}
	}
}

// AddRows adds multiple rows to the table
func (t *Table) AddRows(rows [][]string) {
	for _, row := range rows {
		t.AddRow(row)
	}
}

// SetAlignment sets the alignment for columns
func (t *Table) SetAlignment(alignment []int) {
	t.alignment = alignment
}

// SetSeparator sets the column separator
func (t *Table) SetSeparator(separator string) {
	t.separator = separator
}

// Render renders the table
func (t *Table) Render() error {
	if len(t.headers) == 0 {
		return fmt.Errorf("no headers set")
	}

	// Create tabwriter for formatted output
	w := tabwriter.NewWriter(t.output, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Render headers
	if err := t.renderHeaders(w); err != nil {
		return err
	}

	// Render rows
	for _, row := range t.rows {
		if err := t.renderRow(w, row); err != nil {
			return err
		}
	}

	return nil
}

// renderHeaders renders the table headers
func (t *Table) renderHeaders(w io.Writer) error {
	headerRow := make([]string, len(t.headers))
	for i, header := range t.headers {
		headerRow[i] = t.terminal.Blue(header)
	}
	
	_, err := fmt.Fprintln(w, strings.Join(headerRow, t.separator))
	return err
}

// renderRow renders a table row
func (t *Table) renderRow(w io.Writer, row []string) error {
	// Pad row to match header length
	paddedRow := make([]string, len(t.headers))
	for i := 0; i < len(t.headers); i++ {
		if i < len(row) {
			paddedRow[i] = row[i]
		} else {
			paddedRow[i] = ""
		}
	}
	
	_, err := fmt.Fprintln(w, strings.Join(paddedRow, t.separator))
	return err
}

// GetWidth returns the table width
func (t *Table) GetWidth() int {
	if len(t.widths) == 0 {
		return 0
	}
	
	width := 0
	for _, w := range t.widths {
		width += w
	}
	width += (len(t.widths) - 1) * len(t.separator)
	
	return width
}

// GetHeight returns the table height
func (t *Table) GetHeight() int {
	return len(t.rows) + 1 // +1 for header
}

// Clear clears the table
func (t *Table) Clear() {
	t.headers = nil
	t.rows = nil
	t.widths = nil
	t.alignment = nil
}

// IsEmpty returns whether the table is empty
func (t *Table) IsEmpty() bool {
	return len(t.rows) == 0
}

// GetRowCount returns the number of rows
func (t *Table) GetRowCount() int {
	return len(t.rows)
}

// GetColumnCount returns the number of columns
func (t *Table) GetColumnCount() int {
	return len(t.headers)
}

// SetOutput sets the output writer
func (t *Table) SetOutput(w io.Writer) {
	t.output = w
}

// GetOutput returns the output writer
func (t *Table) GetOutput() io.Writer {
	return t.output
}

// SetTerminal sets the terminal instance
func (t *Table) SetTerminal(term *terminal.Terminal) {
	t.terminal = term
}

// GetTerminal returns the terminal instance
func (t *Table) GetTerminal() *terminal.Terminal {
	return t.terminal
}

// Format formats the table as a string
func (t *Table) Format() string {
	var buf strings.Builder
	tempTable := &Table{
		headers:   t.headers,
		rows:      t.rows,
		output:    &buf,
		terminal:  t.terminal,
		separator: t.separator,
		alignment: t.alignment,
		widths:    t.widths,
	}
	
	if err := tempTable.Render(); err != nil {
		return ""
	}
	
	return buf.String()
}

// Print prints the table to stdout
func (t *Table) Print() error {
	return t.Render()
}

// PrintTo prints the table to a writer
func (t *Table) PrintTo(w io.Writer) error {
	originalOutput := t.output
	t.output = w
	defer func() { t.output = originalOutput }()
	
	return t.Render()
}
