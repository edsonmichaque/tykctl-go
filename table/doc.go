// Package table provides beautiful table formatting and display for CLI applications.
//
// Features:
//   - Table Rendering: Format and display data in tables
//   - Header Support: Configurable table headers with styling
//   - Row Management: Add rows and manage table data
//   - Terminal Integration: Integration with terminal utilities for styling
//   - Customizable: Configurable separators, alignment, and styling
//   - Output Control: Support for different output writers
//
// Example:
//   tbl := table.New()
//   tbl.SetHeaders([]string{"Name", "Age", "City"})
//   tbl.AddRow([]string{"John", "30", "NYC"})
//   tbl.Render()
package table