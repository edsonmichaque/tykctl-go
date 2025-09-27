package table

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/edsonmichaque/tykctl-go/terminal"
	"github.com/olekukonko/tablewriter/tw"
)

// Table represents a table for formatted output
type Table struct {
	headers         []string
	rows            [][]string
	output          io.Writer
	terminal        *terminal.Terminal
	separator       string
	alignment       []int
	widths          []int
	border          bool
	centerSeparator string
	columnSeparator string
	rowSeparator    string
	headerLine      bool
}

// New creates a new table instance
func New() *Table {
	return &Table{
		output:          os.Stdout,
		terminal:        terminal.New(),
		separator:       "\t",
		border:          false, // Default to borderless
		centerSeparator: "+",
		columnSeparator: "|",
		rowSeparator:    "-",
		headerLine:      false,
	}
}

// NewWithWriter creates a new table with a custom writer
func NewWithWriter(w io.Writer) *Table {
	return &Table{
		output:          w,
		terminal:        terminal.New(),
		separator:       "\t",
		border:          false, // Default to borderless
		centerSeparator: "+",
		columnSeparator: "|",
		rowSeparator:    "-",
		headerLine:      false,
	}
}

// SetHeaders sets the table headers
func (t *Table) SetHeaders(headers []string) {
	t.headers = make([]string, len(headers))
	t.widths = make([]int, len(headers))
	for i, header := range headers {
		t.headers[i] = strings.ToUpper(header)
		t.widths[i] = len(t.headers[i])
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

	if t.border {
		return t.renderWithBorders()
	}
	return t.renderWithoutBorders()
}

// renderWithoutBorders renders the table without borders (default behavior)
func (t *Table) renderWithoutBorders() error {
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

// renderWithBorders renders the table with borders
func (t *Table) renderWithBorders() error {
	// Render top border
	if err := t.renderTopBorder(); err != nil {
		return err
	}

	// Render headers with borders
	if err := t.renderHeadersWithBorders(); err != nil {
		return err
	}

	// Render header separator line if enabled
	if t.headerLine {
		if err := t.renderHeaderSeparator(); err != nil {
			return err
		}
	}

	// Render rows with borders
	for _, row := range t.rows {
		if err := t.renderRowWithBorders(row); err != nil {
			return err
		}
	}

	// Render bottom border
	return t.renderBottomBorder()
}

// renderHeaders renders the table headers
func (t *Table) renderHeaders(w io.Writer) error {
	headerRow := make([]string, len(t.headers))
	for i, header := range t.headers {
		headerRow[i] = t.terminal.Blue(strings.ToUpper(header))
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

// Border rendering helper methods

// renderTopBorder renders the top border of the table
func (t *Table) renderTopBorder() error {
	border := t.createHorizontalBorder()
	_, err := fmt.Fprintln(t.output, border)
	return err
}

// renderBottomBorder renders the bottom border of the table
func (t *Table) renderBottomBorder() error {
	border := t.createHorizontalBorder()
	_, err := fmt.Fprintln(t.output, border)
	return err
}

// renderHeaderSeparator renders the line under headers
func (t *Table) renderHeaderSeparator() error {
	border := t.createHorizontalBorder()
	_, err := fmt.Fprintln(t.output, border)
	return err
}

// renderHeadersWithBorders renders headers with column separators
func (t *Table) renderHeadersWithBorders() error {
	headerRow := make([]string, len(t.headers))
	for i, header := range t.headers {
		headerRow[i] = t.terminal.Blue(strings.ToUpper(header))
	}

	line := t.columnSeparator + " " + strings.Join(headerRow, " "+t.columnSeparator+" ") + " " + t.columnSeparator
	_, err := fmt.Fprintln(t.output, line)
	return err
}

// renderRowWithBorders renders a row with column separators
func (t *Table) renderRowWithBorders(row []string) error {
	// Pad row to match header length
	paddedRow := make([]string, len(t.headers))
	for i := 0; i < len(t.headers); i++ {
		if i < len(row) {
			paddedRow[i] = row[i]
		} else {
			paddedRow[i] = ""
		}
	}

	line := t.columnSeparator + " " + strings.Join(paddedRow, " "+t.columnSeparator+" ") + " " + t.columnSeparator
	_, err := fmt.Fprintln(t.output, line)
	return err
}

// createHorizontalBorder creates a horizontal border line
func (t *Table) createHorizontalBorder() string {
	if len(t.widths) == 0 {
		return ""
	}

	parts := make([]string, len(t.widths))
	for i, width := range t.widths {
		parts[i] = strings.Repeat(t.rowSeparator, width+2) // +2 for padding
	}

	return t.centerSeparator + strings.Join(parts, t.centerSeparator) + t.centerSeparator
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
		headers:         t.headers,
		rows:            t.rows,
		output:          &buf,
		terminal:        t.terminal,
		separator:       t.separator,
		alignment:       t.alignment,
		widths:          t.widths,
		border:          t.border,
		centerSeparator: t.centerSeparator,
		columnSeparator: t.columnSeparator,
		rowSeparator:    t.rowSeparator,
		headerLine:      t.headerLine,
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

// Border configuration methods

// SetBorder sets whether the table should have borders
func (t *Table) SetBorder(border bool) {
	t.border = border
}

// SetCenterSeparator sets the center separator character for borders
func (t *Table) SetCenterSeparator(separator string) {
	t.centerSeparator = separator
}

// SetColumnSeparator sets the column separator character for borders
func (t *Table) SetColumnSeparator(separator string) {
	t.columnSeparator = separator
}

// SetRowSeparator sets the row separator character for borders
func (t *Table) SetRowSeparator(separator string) {
	t.rowSeparator = separator
}

// SetHeaderLine sets whether there should be a line under headers
func (t *Table) SetHeaderLine(line bool) {
	t.headerLine = line
}

// SetColumnAlignment sets the column alignment (no-op for borderless tables)
func (t *Table) SetColumnAlignment(alignment []tw.Align) {
	// No-op: this implementation uses fixed left alignment
}

// SetHeaderAlignment sets the header alignment (no-op for borderless tables)
func (t *Table) SetHeaderAlignment(alignment []tw.Align) {
	// No-op: this implementation uses fixed left alignment
}

// SetAutoWrapText sets the auto wrap text flag (no-op for borderless tables)
func (t *Table) SetAutoWrapText(wrap bool) {
	// No-op: this implementation doesn't support text wrapping
}

// SetReflowDuringAutoWrap sets the reflow during auto wrap flag (no-op for borderless tables)
func (t *Table) SetReflowDuringAutoWrap(reflow bool) {
	// No-op: this implementation doesn't support text wrapping
}

// SetTablePadding sets the table padding (no-op for borderless tables)
func (t *Table) SetTablePadding(padding string) {
	// No-op: this implementation uses fixed padding
}

// SetNoWhiteSpace sets the no white space flag (no-op for borderless tables)
func (t *Table) SetNoWhiteSpace(noWhiteSpace bool) {
	// No-op: this implementation uses fixed spacing
}

// SetFooter sets the table footer (no-op for borderless tables)
func (t *Table) SetFooter(footer []string) {
	// No-op: this implementation doesn't support footers
}
