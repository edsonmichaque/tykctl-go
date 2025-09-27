# Table Package

The `table` package provides formatted table output functionality for tykctl applications, offering flexible table rendering with customizable headers, alignment, and styling.

## Features

- **Flexible Tables**: Support for various table formats and layouts
- **Customizable Headers**: Set custom headers with proper alignment
- **Multiple Output Formats**: Support for different output formats
- **Terminal Integration**: Works seamlessly with terminal capabilities
- **Alignment Options**: Configurable column alignment
- **Styling Support**: Customizable table appearance

## Usage

### Basic Table

```go
package main

import (
    "fmt"
    
    "github.com/edsonmichaque/tykctl-go/table"
)

func main() {
    // Create a new table
    t := table.New()
    
    // Set headers
    t.SetHeaders([]string{"Name", "Age", "City"})
    
    // Add rows
    t.AddRow([]string{"Alice", "30", "New York"})
    t.AddRow([]string{"Bob", "25", "London"})
    t.AddRow([]string{"Charlie", "35", "Tokyo"})
    
    // Render table
    t.Render()
}
```

### Table with Custom Writer

```go
func tableWithCustomWriter() {
    // Create table with custom writer
    t := table.NewWithWriter(os.Stderr)
    
    // Set headers
    t.SetHeaders([]string{"Error", "Count", "Severity"})
    
    // Add rows
    t.AddRow([]string{"Database Error", "5", "High"})
    t.AddRow([]string{"Network Timeout", "12", "Medium"})
    t.AddRow([]string{"Validation Error", "3", "Low"})
    
    // Render to stderr
    t.Render()
}
```

### Table with Alignment

```go
func tableWithAlignment() {
    t := table.New()
    
    // Set headers
    t.SetHeaders([]string{"Product", "Price", "Stock"})
    
    // Set alignment (0=left, 1=center, 2=right)
    t.SetAlignment([]int{table.AlignLeft, table.AlignRight, table.AlignRight})
    
    // Add rows
    t.AddRow([]string{"Laptop", "$999.99", "15"})
    t.AddRow([]string{"Mouse", "$29.99", "150"})
    t.AddRow([]string{"Keyboard", "$79.99", "75"})
    
    t.Render()
}
```

## Advanced Usage

### Dynamic Table Building

```go
func dynamicTable() {
    t := table.New()
    
    // Set headers
    headers := []string{"ID", "Name", "Status", "Created"}
    t.SetHeaders(headers)
    
    // Add rows dynamically
    data := [][]string{
        {"1", "Task 1", "Completed", "2023-01-01"},
        {"2", "Task 2", "In Progress", "2023-01-02"},
        {"3", "Task 3", "Pending", "2023-01-03"},
    }
    
    for _, row := range data {
        t.AddRow(row)
    }
    
    t.Render()
}
```

### Table with Custom Separator

```go
func tableWithSeparator() {
    t := table.New()
    
    // Set custom separator
    t.SetSeparator(" | ")
    
    // Set headers
    t.SetHeaders([]string{"Command", "Description", "Usage"})
    
    // Add rows
    t.AddRow([]string{"install", "Install extension", "tykctl install <repo>"})
    t.AddRow([]string{"list", "List extensions", "tykctl list"})
    t.AddRow([]string{"remove", "Remove extension", "tykctl remove <name>"})
    
    t.Render()
}
```

### Table with Width Constraints

```go
func tableWithWidth() {
    t := table.New()
    
    // Set headers
    t.SetHeaders([]string{"Long Column Name", "Short", "Very Long Column Name"})
    
    // Set column widths
    t.SetWidths([]int{20, 10, 25})
    
    // Add rows
    t.AddRow([]string{"This is a long value", "Short", "This is another very long value"})
    t.AddRow([]string{"Another long value", "Tiny", "Yet another long value here"})
    
    t.Render()
}
```

## Integration Examples

### With Extension List

```go
func listExtensions() error {
    installer := extension.NewInstaller("/tmp/tykctl-config")
    ctx := context.Background()
    
    // Get installed extensions
    installed, err := installer.ListInstalledExtensions(ctx)
    if err != nil {
        return err
    }
    
    // Create table
    t := table.New()
    t.SetHeaders([]string{"Name", "Version", "Repository", "Installed"})
    
    // Add extension data
    for _, ext := range installed {
        t.AddRow([]string{
            ext.Name,
            ext.Version,
            ext.Repository,
            ext.InstalledAt.Format("2006-01-02"),
        })
    }
    
    // Render table
    t.Render()
    return nil
}
```

### With API Response Data

```go
func displayAPIResponse() error {
    client := httpclient.NewWithBaseURL("https://api.github.com")
    ctx := context.Background()
    
    // Get user data
    resp, err := client.Get(ctx, "/users/octocat")
    if err != nil {
        return err
    }
    
    if !resp.IsSuccess() {
        return fmt.Errorf("API request failed: %d", resp.StatusCode)
    }
    
    // Parse response
    var user map[string]interface{}
    err = resp.UnmarshalJSON(&user)
    if err != nil {
        return err
    }
    
    // Create table
    t := table.New()
    t.SetHeaders([]string{"Property", "Value"})
    
    // Add user data
    t.AddRow([]string{"Name", fmt.Sprintf("%v", user["name"])})
    t.AddRow([]string{"Login", fmt.Sprintf("%v", user["login"])})
    t.AddRow([]string{"Followers", fmt.Sprintf("%v", user["followers"])})
    t.AddRow([]string{"Public Repos", fmt.Sprintf("%v", user["public_repos"])})
    
    t.Render()
    return nil
}
```

### With Configuration Display

```go
func displayConfiguration() error {
    // Load configuration
    config, err := loadConfig()
    if err != nil {
        return err
    }
    
    // Create table
    t := table.New()
    t.SetHeaders([]string{"Setting", "Value", "Description"})
    
    // Add configuration data
    t.AddRow([]string{"Server Host", config.Server.Host, "Server hostname"})
    t.AddRow([]string{"Server Port", fmt.Sprintf("%d", config.Server.Port), "Server port"})
    t.AddRow([]string{"Database URL", config.Database.URL, "Database connection string"})
    t.AddRow([]string{"Debug Mode", fmt.Sprintf("%t", config.Debug), "Enable debug logging"})
    
    t.Render()
    return nil
}
```

## Customization Options

### Custom Styling

```go
func customStyling() {
    t := table.New()
    
    // Set headers
    t.SetHeaders([]string{"Status", "Message", "Time"})
    
    // Custom styling
    t.SetSeparator(" | ")
    t.SetAlignment([]int{table.AlignCenter, table.AlignLeft, table.AlignRight})
    
    // Add rows with different statuses
    t.AddRow([]string{"✓", "Operation completed successfully", "10:30:15"})
    t.AddRow([]string{"⚠", "Warning: High memory usage", "10:31:22"})
    t.AddRow([]string{"✗", "Error: Connection failed", "10:32:45"})
    
    t.Render()
}
```

### Table with Colors

```go
func tableWithColors() {
    t := table.New()
    
    // Set headers
    t.SetHeaders([]string{"Level", "Message", "Count"})
    
    // Add colored rows
    t.AddRow([]string{"ERROR", "Database connection failed", "5"})
    t.AddRow([]string{"WARN", "High memory usage", "12"})
    t.AddRow([]string{"INFO", "User logged in", "150"})
    t.AddRow([]string{"DEBUG", "Cache hit", "1200"})
    
    t.Render()
}
```

## Error Handling

### Table with Error Data

```go
func displayErrors() error {
    // Get error data
    errors, err := getErrorData()
    if err != nil {
        return err
    }
    
    if len(errors) == 0 {
        fmt.Println("No errors found")
        return nil
    }
    
    // Create table
    t := table.New()
    t.SetHeaders([]string{"Time", "Level", "Message", "Source"})
    
    // Add error data
    for _, err := range errors {
        t.AddRow([]string{
            err.Timestamp.Format("15:04:05"),
            err.Level,
            err.Message,
            err.Source,
        })
    }
    
    t.Render()
    return nil
}
```

## Best Practices

- **Clear Headers**: Use descriptive column headers
- **Consistent Formatting**: Maintain consistent data formatting across rows
- **Appropriate Alignment**: Use appropriate alignment for different data types
- **Width Management**: Set appropriate column widths for readability
- **Error Handling**: Handle empty data gracefully
- **Performance**: Avoid rendering very large tables at once

## Dependencies

- `github.com/edsonmichaque/tykctl-go/terminal` - Terminal capabilities
- `github.com/olekukonko/tablewriter/tw` - Table formatting library