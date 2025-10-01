# Templates System

The TykCtl-Go templates system provides a powerful way to generate resources from predefined templates with variable substitution.

## Overview

Templates allow you to create reusable resource definitions that can be customized with variables. This is particularly useful for:

- **API Definitions**: Standard API configurations with customizable parameters
- **Policy Templates**: Reusable policy configurations
- **User Templates**: Standard user configurations
- **Environment Configurations**: Environment-specific settings

## Template Structure

### Basic Template

```yaml
# api-template.yaml
name: "{{.name}}"
slug: "{{.slug}}"
version: "{{.version | default "1.0.0"}}"
description: "{{.description}}"
api_definition:
  name: "{{.name}}"
  slug: "{{.slug}}"
  version: "{{.version | default "1.0.0"}}"
  proxy:
    listen_path: "/{{.slug}}"
    target_url: "{{.target_url}}"
    strip_listen_path: true
  auth:
    auth_header_name: "Authorization"
  version_data:
    versions:
      "{{.version | default "1.0.0"}}":
        name: "{{.version | default "1.0.0"}}"
```

### Template Variables

Templates support Go template syntax with additional functions:

- **`{{.variable}}`**: Simple variable substitution
- **`{{.variable | default "value"}}`**: Default value if variable is empty
- **`{{.variable | upper}}`**: Convert to uppercase
- **`{{.variable | lower}}`**: Convert to lowercase
- **`{{.variable | title}}`**: Title case
- **`{{.variable | trim}}`**: Trim whitespace

## Usage Examples

### Creating Resources from Templates

```go
package main

import (
    "context"
    "fmt"
    "github.com/edsonmichaque/tykctl-go/template"
)

func main() {
    ctx := context.Background()
    
    // Load template
    tmpl, err := template.LoadTemplate("api-template.yaml")
    if err != nil {
        panic(err)
    }
    
    // Define variables
    variables := map[string]interface{}{
        "name":        "My API",
        "slug":        "my-api",
        "version":     "2.0.0",
        "description": "A sample API",
        "target_url":  "https://api.example.com",
    }
    
    // Generate resource
    resource, err := tmpl.Render(variables)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(resource)
}
```

### Template Discovery

```go
// Discover templates
templates, err := template.DiscoverTemplates(ctx, "portal")
if err != nil {
    panic(err)
}

for _, tmpl := range templates {
    fmt.Printf("Template: %s (%s)\n", tmpl.Name, tmpl.Description)
}
```

## Template Management

### Template Installation

```bash
# Install template from file
tykctl-portal template install /path/to/template.yaml

# Install template from directory
tykctl-portal template install /path/to/templates/

# Create template from existing resource
tykctl-portal template create api-template.yaml --from-resource api-123
```

### Template Discovery

Templates are discovered from multiple locations:

1. **Extension-specific templates**: `~/.config/tykctl/<extension>/templates/`
2. **System templates**: `/usr/local/share/tykctl/<extension>/templates/`
3. **Development templates**: `./templates/`
4. **Custom paths**: Configured via environment variables

## Environment Variables

### Template Configuration

```bash
# Template discovery paths
TYKCTL_PORTAL_TEMPLATE_DIR="/path/to/templates"
TYKCTL_PORTAL_TEMPLATE_DATA_DIR="/path/to/template/data"

# Template settings
TYKCTL_TEMPLATE_CACHE_ENABLED=true
TYKCTL_TEMPLATE_CACHE_TTL="1h"
```

## Advanced Features

### Template Inheritance

```yaml
# base-api-template.yaml
name: "{{.name}}"
version: "{{.version}}"
proxy:
  listen_path: "/{{.slug}}"
  target_url: "{{.target_url}}"

# extended-api-template.yaml
{{template "base-api-template.yaml" .}}
auth:
  auth_header_name: "{{.auth_header | default "Authorization"}}"
rate_limit:
  requests_per_minute: {{.rate_limit | default 100}}
```

### Conditional Logic

```yaml
name: "{{.name}}"
{{if .auth_enabled}}
auth:
  auth_header_name: "Authorization"
{{end}}
{{if .rate_limit}}
rate_limit:
  requests_per_minute: {{.rate_limit}}
{{end}}
```

### Template Functions

```yaml
# Custom functions available in templates
name: "{{.name | title}}"
slug: "{{.name | lower | replace " " "-"}}"
version: "{{.version | default "1.0.0"}}"
url: "{{.protocol | default "https"}}://{{.host}}/{{.path}}"
```

## Best Practices

### 1. Template Design

- **Use descriptive names**: `api-template.yaml`, `user-template.yaml`
- **Include documentation**: Add comments explaining variables
- **Provide defaults**: Use default values for optional parameters
- **Validate inputs**: Include validation in templates when possible

### 2. Variable Naming

```yaml
# Good: Descriptive variable names
name: "{{.api_name}}"
target_url: "{{.backend_url}}"
auth_header: "{{.auth_header_name}}"

# Avoid: Generic or unclear names
name: "{{.name}}"
url: "{{.url}}"
header: "{{.header}}"
```

### 3. Template Organization

```
templates/
├── api/
│   ├── rest-api-template.yaml
│   ├── graphql-api-template.yaml
│   └── websocket-api-template.yaml
├── user/
│   ├── admin-user-template.yaml
│   └── regular-user-template.yaml
└── policy/
    ├── rate-limit-policy.yaml
    └── auth-policy.yaml
```

## Troubleshooting

### Common Issues

**Template not found:**
- Check template discovery paths
- Verify template file exists
- Check file permissions

**Variable substitution errors:**
- Verify variable names match template
- Check for typos in variable names
- Ensure required variables are provided

**Template syntax errors:**
- Validate YAML syntax
- Check Go template syntax
- Verify function usage

### Debug Mode

```bash
# Enable template debug mode
export TYKCTL_TEMPLATE_DEBUG=true

# Verbose template processing
tykctl-portal template render api-template.yaml --verbose
```

## Examples

### API Template

```yaml
# rest-api-template.yaml
name: "{{.name}}"
slug: "{{.slug}}"
version: "{{.version | default "1.0.0"}}"
description: "{{.description}}"
api_definition:
  name: "{{.name}}"
  slug: "{{.slug}}"
  version: "{{.version | default "1.0.0"}}"
  proxy:
    listen_path: "/{{.slug}}"
    target_url: "{{.target_url}}"
    strip_listen_path: true
  auth:
    auth_header_name: "{{.auth_header | default "Authorization"}}"
  version_data:
    versions:
      "{{.version | default "1.0.0"}}":
        name: "{{.version | default "1.0.0"}}"
        use_extended_paths: true
        paths:
          ignored: []
          white_list: []
          black_list: []
```

### User Template

```yaml
# admin-user-template.yaml
first_name: "{{.first_name}}"
last_name: "{{.last_name}}"
email: "{{.email}}"
password: "{{.password}}"
user_permissions:
  admin: true
  developer: true
  analytics: true
```

### Policy Template

```yaml
# rate-limit-policy.yaml
name: "{{.name}}"
slug: "{{.slug}}"
rate: {{.rate_limit | default 100}}
per: {{.per_minute | default 60}}
quota_max: {{.quota_max | default 1000}}
quota_renewal_rate: {{.quota_renewal | default 3600}}
access_rights:
  "{{.api_id}}":
    api_name: "{{.api_name}}"
    api_id: "{{.api_id}}"
    versions:
      - "{{.version | default "Default"}}"
```

## Resources

- [Configuration Guide](../config/README.md)
- [Extension Framework](../extension/README.md)
- [Getting Started Guide](../guides/getting-started.md)