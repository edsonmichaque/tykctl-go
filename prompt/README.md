# Prompt Package

The `prompt` package provides interactive command-line prompting functionality for tykctl applications, offering various input types including strings, numbers, booleans, and selections.

## Features

- **Interactive Prompts**: Various prompt types for user input
- **Input Validation**: Built-in validation for different input types
- **Default Values**: Support for default values in prompts
- **Confirmation Prompts**: Yes/no confirmation prompts
- **Selection Prompts**: Choose from multiple options
- **Terminal Integration**: Works seamlessly with terminal capabilities

## Usage

### Basic String Prompt

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/edsonmichaque/tykctl-go/prompt"
)

func main() {
    // Create a new prompt instance
    p := prompt.New()
    
    // Ask for string input
    name, err := p.AskString("What is your name?")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Hello, %s!\n", name)
}
```

### String Prompt with Default

```go
func promptWithDefault() {
    p := prompt.New()
    
    // Ask for input with default value
    username, err := p.AskStringWithDefault("Enter username", "admin")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Username: %s\n", username)
}
```

### Number Prompts

```go
func numberPrompts() {
    p := prompt.New()
    
    // Ask for integer
    age, err := p.AskInt("What is your age?")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Age: %d\n", age)
    
    // Ask for integer with default
    port, err := p.AskIntWithDefault("Enter port number", 8080)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Port: %d\n", port)
    
    // Ask for float
    price, err := p.AskFloat("Enter price")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Price: %.2f\n", price)
}
```

### Boolean Prompts

```go
func booleanPrompts() {
    p := prompt.New()
    
    // Ask for yes/no
    enabled, err := p.AskBool("Enable feature?")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Feature enabled: %t\n", enabled)
    
    // Ask for yes/no with default
    debug, err := p.AskBoolWithDefault("Enable debug mode?", false)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Debug mode: %t\n", debug)
}
```

## Advanced Usage

### Selection Prompts

```go
func selectionPrompts() {
    p := prompt.New()
    
    // Single selection
    options := []string{"Option 1", "Option 2", "Option 3"}
    choice, err := p.AskSelection("Choose an option:", options)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Selected: %s\n", choice)
    
    // Multiple selection
    multiChoice, err := p.AskMultiSelection("Choose multiple options:", options)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Selected: %v\n", multiChoice)
}
```

### Password Prompts

```go
func passwordPrompts() {
    p := prompt.New()
    
    // Ask for password (hidden input)
    password, err := p.AskPassword("Enter password:")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Password length: %d\n", len(password))
    
    // Confirm password
    confirmPassword, err := p.AskPassword("Confirm password:")
    if err != nil {
        log.Fatal(err)
    }
    
    if password != confirmPassword {
        fmt.Println("Passwords do not match!")
        return
    }
    
    fmt.Println("Passwords match!")
}
```

### Validation Prompts

```go
func validationPrompts() {
    p := prompt.New()
    
    // Email validation
    email, err := p.AskStringWithValidation("Enter email:", func(input string) error {
        if !strings.Contains(input, "@") {
            return fmt.Errorf("invalid email format")
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Email: %s\n", email)
    
    // Number range validation
    age, err := p.AskIntWithValidation("Enter age:", func(input int) error {
        if input < 0 || input > 150 {
            return fmt.Errorf("age must be between 0 and 150")
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Age: %d\n", age)
}
```

## Integration Examples

### With Configuration Setup

```go
func setupConfiguration() error {
    p := prompt.New()
    
    fmt.Println("=== Configuration Setup ===")
    
    // Server configuration
    host, err := p.AskStringWithDefault("Enter server host", "localhost")
    if err != nil {
        return err
    }
    
    port, err := p.AskIntWithDefault("Enter server port", 8080)
    if err != nil {
        return err
    }
    
    // Database configuration
    dbHost, err := p.AskStringWithDefault("Enter database host", "localhost")
    if err != nil {
        return err
    }
    
    dbPort, err := p.AskIntWithDefault("Enter database port", 5432)
    if err != nil {
        return err
    }
    
    // Authentication
    useAuth, err := p.AskBoolWithDefault("Enable authentication?", true)
    if err != nil {
        return err
    }
    
    var username, password string
    if useAuth {
        username, err = p.AskString("Enter username:")
        if err != nil {
            return err
        }
        
        password, err = p.AskPassword("Enter password:")
        if err != nil {
            return err
        }
    }
    
    // Save configuration
    config := map[string]interface{}{
        "server": map[string]interface{}{
            "host": host,
            "port": port,
        },
        "database": map[string]interface{}{
            "host": dbHost,
            "port": dbPort,
        },
        "auth": map[string]interface{}{
            "enabled":  useAuth,
            "username": username,
            "password": password,
        },
    }
    
    fmt.Printf("Configuration saved: %+v\n", config)
    return nil
}
```

### With Extension Installation

```go
func installExtensionWithPrompts() error {
    p := prompt.New()
    
    fmt.Println("=== Extension Installation ===")
    
    // Get extension details
    owner, err := p.AskString("Enter GitHub owner:")
    if err != nil {
        return err
    }
    
    repo, err := p.AskString("Enter repository name:")
    if err != nil {
        return err
    }
    
    // Confirm installation
    confirm, err := p.AskBoolWithDefault(fmt.Sprintf("Install %s/%s?", owner, repo), true)
    if err != nil {
        return err
    }
    
    if !confirm {
        fmt.Println("Installation cancelled")
        return nil
    }
    
    // Install extension
    installer := extension.NewInstaller("/tmp/tykctl-config")
    ctx := context.Background()
    
    err = installer.InstallExtension(ctx, owner, repo)
    if err != nil {
        return err
    }
    
    fmt.Println("Extension installed successfully!")
    return nil
}
```

### With User Registration

```go
func userRegistration() error {
    p := prompt.New()
    
    fmt.Println("=== User Registration ===")
    
    // Get user information
    username, err := p.AskStringWithValidation("Enter username:", func(input string) error {
        if len(input) < 3 {
            return fmt.Errorf("username must be at least 3 characters")
        }
        return nil
    })
    if err != nil {
        return err
    }
    
    email, err := p.AskStringWithValidation("Enter email:", func(input string) error {
        if !strings.Contains(input, "@") {
            return fmt.Errorf("invalid email format")
        }
        return nil
    })
    if err != nil {
        return err
    }
    
    password, err := p.AskPassword("Enter password:")
    if err != nil {
        return err
    }
    
    confirmPassword, err := p.AskPassword("Confirm password:")
    if err != nil {
        return err
    }
    
    if password != confirmPassword {
        return fmt.Errorf("passwords do not match")
    }
    
    // Role selection
    roles := []string{"user", "admin", "moderator"}
    role, err := p.AskSelection("Select role:", roles)
    if err != nil {
        return err
    }
    
    // Terms and conditions
    acceptTerms, err := p.AskBoolWithDefault("Accept terms and conditions?", false)
    if err != nil {
        return err
    }
    
    if !acceptTerms {
        return fmt.Errorf("terms and conditions must be accepted")
    }
    
    // Create user
    user := map[string]interface{}{
        "username": username,
        "email":    email,
        "role":     role,
    }
    
    fmt.Printf("User created: %+v\n", user)
    return nil
}
```

## Custom Prompt Types

### Custom Validation

```go
func customValidation() {
    p := prompt.New()
    
    // URL validation
    url, err := p.AskStringWithValidation("Enter URL:", func(input string) error {
        if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
            return fmt.Errorf("URL must start with http:// or https://")
        }
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("URL: %s\n", url)
}
```

### Custom Selection with Descriptions

```go
func customSelection() {
    p := prompt.New()
    
    // Custom selection with descriptions
    options := []string{
        "Basic Plan - $9/month",
        "Pro Plan - $29/month", 
        "Enterprise Plan - $99/month",
    }
    
    choice, err := p.AskSelection("Choose a plan:", options)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Selected plan: %s\n", choice)
}
```

## Error Handling

### Robust Error Handling

```go
func robustPrompting() error {
    p := prompt.New()
    
    // Retry on validation failure
    maxRetries := 3
    for attempt := 1; attempt <= maxRetries; attempt++ {
        username, err := p.AskStringWithValidation("Enter username:", func(input string) error {
            if len(input) < 3 {
                return fmt.Errorf("username must be at least 3 characters")
            }
            if strings.Contains(input, " ") {
                return fmt.Errorf("username cannot contain spaces")
            }
            return nil
        })
        
        if err == nil {
            fmt.Printf("Username: %s\n", username)
            return nil
        }
        
        if attempt < maxRetries {
            fmt.Printf("Invalid input (attempt %d/%d): %v\n", attempt, maxRetries, err)
        }
    }
    
    return fmt.Errorf("failed to get valid input after %d attempts", maxRetries)
}
```

## Best Practices

- **Clear Prompts**: Use clear, descriptive prompt messages
- **Default Values**: Provide sensible defaults when possible
- **Validation**: Validate input to prevent errors
- **Error Handling**: Handle errors gracefully with retry options
- **User Experience**: Provide helpful error messages
- **Confirmation**: Ask for confirmation for destructive operations

## Dependencies

- `github.com/edsonmichaque/tykctl-go/terminal` - Terminal capabilities