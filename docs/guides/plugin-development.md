# Plugin Development Guide

This guide covers how to develop plugins for TykCtl extensions using the TykCtl-Go plugin system.

## 🎯 Plugin Overview

Plugins are third-party executables that extend TykCtl extension functionality. They can be written in any language and provide custom logic, integrations, and workflows.

### Plugin Characteristics

- **Extension-Specific**: Belong to a specific extension (Portal, Dashboard, Gateway)
- **Custom Logic**: Provide specialized functionality not available in core extensions
- **Third-Party**: Developed by users, organizations, or the community
- **Optional**: Extensions work without plugins; plugins enhance functionality

## 📋 Plugin Requirements

### Naming Convention

Plugins must follow the naming convention: `tykctl-<extension>-<name>`

**Examples:**
- `tykctl-portal-deploy`
- `tykctl-dashboard-backup`
- `tykctl-gateway-monitor`

### Executable Requirements

- Must be executable by the operating system
- Should handle command-line arguments
- Should return appropriate exit codes
- Should be cross-platform compatible when possible

## 🛠️ Development Process

### 1. Choose Your Language

Plugins can be written in any language that can produce executables:

- **Shell Scripts** (bash, zsh, fish)
- **Python** scripts
- **Go** programs
- **Node.js** applications
- **Rust** programs
- **Any compiled language**

### 2. Create Plugin Structure

#### Shell Script Plugin

```bash
#!/bin/bash
# tykctl-portal-deploy

set -euo pipefail

PLUGIN_NAME="deploy"
PLUGIN_VERSION="1.0.0"

# Logging function
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$PLUGIN_NAME] $*" >&2
}

# Help function
show_help() {
    cat << EOF
Usage: $0 [OPTIONS] [COMMAND]

Commands:
    production    Deploy to production environment
    staging       Deploy to staging environment
    rollback      Rollback last deployment

Options:
    --dry-run     Show what would be deployed without executing
    --force       Force deployment even if checks fail
    --help        Show this help message

Environment Variables:
    TYK_PORTAL_URL      Portal API URL
    TYK_PORTAL_TOKEN    Portal API token
    TYKCTL_PORTAL_CONTEXT  Current context
EOF
}

# Main function
main() {
    local command="${1:-help}"
    
    case "$command" in
        "production")
            deploy_production "$@"
            ;;
        "staging")
            deploy_staging "$@"
            ;;
        "rollback")
            rollback_deployment "$@"
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            log "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Deployment functions
deploy_production() {
    log "Deploying to production environment..."
    
    # Check environment variables
    if [[ -z "${TYK_PORTAL_URL:-}" ]]; then
        log "Error: TYK_PORTAL_URL not set"
        exit 1
    fi
    
    if [[ -z "${TYK_PORTAL_TOKEN:-}" ]]; then
        log "Error: TYK_PORTAL_TOKEN not set"
        exit 1
    fi
    
    # Perform deployment
    log "Portal URL: $TYK_PORTAL_URL"
    log "Context: ${TYKCTL_PORTAL_CONTEXT:-default}"
    
    # Your deployment logic here
    echo "Deployment completed successfully!"
}

deploy_staging() {
    log "Deploying to staging environment..."
    # Staging deployment logic
}

rollback_deployment() {
    log "Rolling back deployment..."
    # Rollback logic
}

# Execute main function
main "$@"
```

#### Python Plugin

```python
#!/usr/bin/env python3
"""
tykctl-portal-backup - Backup plugin for TykCtl Portal
"""

import argparse
import json
import os
import sys
import time
from datetime import datetime
from typing import Dict, List, Optional

class PortalBackupPlugin:
    def __init__(self):
        self.plugin_name = "backup"
        self.version = "1.0.0"
        
    def log(self, message: str) -> None:
        """Log a message with timestamp."""
        timestamp = datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        print(f"[{timestamp}] [{self.plugin_name}] {message}", file=sys.stderr)
        
    def get_config(self) -> Dict[str, str]:
        """Get configuration from environment variables."""
        config = {
            'portal_url': os.getenv('TYK_PORTAL_URL', ''),
            'portal_token': os.getenv('TYK_PORTAL_TOKEN', ''),
            'context': os.getenv('TYKCTL_PORTAL_CONTEXT', 'default'),
            'backup_dir': os.getenv('TYKCTL_PORTAL_BACKUP_DIR', './backups'),
        }
        
        # Validate required config
        if not config['portal_url']:
            self.log("Error: TYK_PORTAL_URL not set")
            sys.exit(1)
            
        if not config['portal_token']:
            self.log("Error: TYK_PORTAL_TOKEN not set")
            sys.exit(1)
            
        return config
        
    def create_backup(self, config: Dict[str, str], dry_run: bool = False) -> None:
        """Create a backup of Portal data."""
        self.log("Starting backup process...")
        
        if dry_run:
            self.log("DRY RUN: Would create backup")
            return
            
        # Create backup directory
        backup_dir = config['backup_dir']
        os.makedirs(backup_dir, exist_ok=True)
        
        # Generate backup filename
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        backup_file = f"{backup_dir}/portal_backup_{timestamp}.json"
        
        # Perform backup
        self.log(f"Creating backup: {backup_file}")
        
        # Your backup logic here
        backup_data = {
            'timestamp': timestamp,
            'context': config['context'],
            'portal_url': config['portal_url'],
            'data': {}  # Your actual backup data
        }
        
        with open(backup_file, 'w') as f:
            json.dump(backup_data, f, indent=2)
            
        self.log(f"Backup completed: {backup_file}")
        
    def restore_backup(self, backup_file: str, dry_run: bool = False) -> None:
        """Restore from a backup file."""
        self.log(f"Restoring from backup: {backup_file}")
        
        if not os.path.exists(backup_file):
            self.log(f"Error: Backup file not found: {backup_file}")
            sys.exit(1)
            
        if dry_run:
            self.log("DRY RUN: Would restore backup")
            return
            
        # Your restore logic here
        self.log("Restore completed successfully")
        
    def list_backups(self, config: Dict[str, str]) -> None:
        """List available backups."""
        backup_dir = config['backup_dir']
        
        if not os.path.exists(backup_dir):
            self.log("No backup directory found")
            return
            
        backups = [f for f in os.listdir(backup_dir) if f.startswith('portal_backup_')]
        
        if not backups:
            self.log("No backups found")
            return
            
        self.log("Available backups:")
        for backup in sorted(backups):
            backup_path = os.path.join(backup_dir, backup)
            stat = os.stat(backup_path)
            size = stat.st_size
            mtime = datetime.fromtimestamp(stat.st_mtime)
            print(f"  {backup} ({size} bytes, {mtime})")
            
    def main(self) -> None:
        """Main entry point."""
        parser = argparse.ArgumentParser(description='TykCtl Portal Backup Plugin')
        parser.add_argument('command', choices=['create', 'restore', 'list'], 
                          help='Backup command to execute')
        parser.add_argument('--backup-file', help='Backup file for restore command')
        parser.add_argument('--dry-run', action='store_true', 
                          help='Show what would be done without executing')
        
        args = parser.parse_args()
        
        config = self.get_config()
        
        if args.command == 'create':
            self.create_backup(config, args.dry_run)
        elif args.command == 'restore':
            if not args.backup_file:
                self.log("Error: --backup-file required for restore command")
                sys.exit(1)
            self.restore_backup(args.backup_file, args.dry_run)
        elif args.command == 'list':
            self.list_backups(config)

if __name__ == '__main__':
    plugin = PortalBackupPlugin()
    plugin.main()
```

#### Go Plugin

```go
package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/edsonmichaque/tykctl-go/api"
)

// Plugin represents the monitor plugin
type Plugin struct {
    name    string
    version string
}

// NewPlugin creates a new monitor plugin
func NewPlugin() *Plugin {
    return &Plugin{
        name:    "monitor",
        version: "1.0.0",
    }
}

// Config holds plugin configuration
type Config struct {
    PortalURL   string `json:"portal_url"`
    PortalToken string `json:"portal_token"`
    Context     string `json:"context"`
    Interval    int    `json:"interval"`
}

// GetConfig retrieves configuration from environment
func (p *Plugin) GetConfig() (*Config, error) {
    config := &Config{
        PortalURL:   os.Getenv("TYK_PORTAL_URL"),
        PortalToken: os.Getenv("TYK_PORTAL_TOKEN"),
        Context:     os.Getenv("TYKCTL_PORTAL_CONTEXT"),
        Interval:    30, // default 30 seconds
    }
    
    if config.PortalURL == "" {
        return nil, fmt.Errorf("TYK_PORTAL_URL not set")
    }
    
    if config.PortalToken == "" {
        return nil, fmt.Errorf("TYK_PORTAL_TOKEN not set")
    }
    
    if config.Context == "" {
        config.Context = "default"
    }
    
    return config, nil
}

// Monitor performs monitoring tasks
func (p *Plugin) Monitor(ctx context.Context, config *Config) error {
    log.Printf("[%s] Starting monitoring...", p.name)
    
    // Create API client
    client := api.NewClient(config.PortalURL, config.PortalToken)
    
    ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            log.Printf("[%s] Monitoring stopped", p.name)
            return ctx.Err()
        case <-ticker.C:
            if err := p.checkHealth(client); err != nil {
                log.Printf("[%s] Health check failed: %v", p.name, err)
            }
        }
    }
}

// checkHealth performs a health check
func (p *Plugin) checkHealth(client *api.Client) error {
    // Your health check logic here
    log.Printf("[%s] Performing health check...", p.name)
    
    // Example: Check API status
    status, err := client.GetStatus()
    if err != nil {
        return fmt.Errorf("failed to get status: %w", err)
    }
    
    log.Printf("[%s] API Status: %s", p.name, status)
    return nil
}

func main() {
    var (
        interval = flag.Int("interval", 30, "Monitoring interval in seconds")
        help     = flag.Bool("help", false, "Show help")
    )
    flag.Parse()
    
    if *help {
        fmt.Println("TykCtl Portal Monitor Plugin")
        fmt.Println("Usage: tykctl-portal-monitor [options]")
        fmt.Println("Options:")
        flag.PrintDefaults()
        return
    }
    
    plugin := NewPlugin()
    
    config, err := plugin.GetConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    config.Interval = *interval
    
    ctx := context.Background()
    if err := plugin.Monitor(ctx, config); err != nil {
        log.Fatal(err)
    }
}
```

### 3. Make Plugin Executable

```bash
# Shell scripts
chmod +x tykctl-portal-deploy

# Python scripts
chmod +x tykctl-portal-backup.py

# Go programs
go build -o tykctl-portal-monitor main.go
chmod +x tykctl-portal-monitor
```

## 🔧 Plugin Installation

### Manual Installation

```bash
# Copy plugin to plugin directory
cp tykctl-portal-deploy /path/to/plugin/dir/

# Or use the plugin install command
tykctl-portal plugin install tykctl-portal-deploy
```

### Automated Installation

```bash
# Install from file
tykctl-portal plugin install /path/to/plugin-file --name=my-plugin

# Install from directory
tykctl-portal plugin install /path/to/plugin/dir

# Create plugin template
tykctl-portal plugin install my-plugin  # Creates template
```

## 🌍 Environment Variables

Plugins receive comprehensive environment variables:

### Core Plugin Variables

```bash
TYKCTL_PLUGIN_NAME=deploy                    # Plugin name
TYKCTL_PLUGIN_PATH=/path/to/plugin           # Plugin path
TYKCTL_PLUGIN_EXTENSION=portal               # Extension name
TYKCTL_PLUGIN_DIR=/path/to/plugin/dir        # Plugin directory
```

### Extension-Specific Variables

```bash
TYKCTL_PORTAL_CONFIG_DIR=/path/to/config     # Config directory
TYKCTL_PORTAL_PLUGIN_DIR=/path/to/plugins    # Plugin directory
TYKCTL_PORTAL_CONTEXT=production             # Current context
TYKCTL_PORTAL_DEBUG=true                     # Debug mode
TYKCTL_PORTAL_VERBOSE=true                   # Verbose mode
```

### API Configuration

```bash
TYK_PORTAL_URL=https://portal.example.com   # Portal URL
TYK_PORTAL_TOKEN=secret-token                # Portal token
TYK_DASHBOARD_URL=https://dashboard.example.com  # Dashboard URL
TYK_DASHBOARD_TOKEN=dashboard-token          # Dashboard token
```

## 🧪 Testing Plugins

### Unit Testing

```bash
#!/bin/bash
# test-plugin.sh

# Test plugin with mock environment
export TYK_PORTAL_URL="https://test.example.com"
export TYK_PORTAL_TOKEN="test-token"
export TYKCTL_PORTAL_CONTEXT="test"

# Test help command
./tykctl-portal-deploy --help

# Test dry run
./tykctl-portal-deploy production --dry-run

# Test error handling
unset TYK_PORTAL_URL
./tykctl-portal-deploy production  # Should fail gracefully
```

### Integration Testing

```go
func TestPluginExecution(t *testing.T) {
    // Create test plugin
    testPlugin := createTestPlugin(t)
    defer os.Remove(testPlugin)
    
    // Set up test environment
    os.Setenv("TYK_PORTAL_URL", "https://test.example.com")
    os.Setenv("TYK_PORTAL_TOKEN", "test-token")
    defer os.Unsetenv("TYK_PORTAL_URL")
    defer os.Unsetenv("TYK_PORTAL_TOKEN")
    
    // Test execution
    manager := plugin.NewManager("portal", testConfigProvider{})
    err := manager.Execute(context.Background(), testPlugin, []string{"--help"})
    
    assert.NoError(t, err)
}
```

## 📦 Distribution

### Plugin Packaging

Create a simple package structure:

```
my-plugin/
├── README.md
├── install.sh
├── tykctl-portal-my-plugin
└── config/
    └── plugin.yaml
```

### Installation Script

```bash
#!/bin/bash
# install.sh

PLUGIN_NAME="my-plugin"
PLUGIN_FILE="tykctl-portal-my-plugin"
INSTALL_DIR="${TYKCTL_PORTAL_PLUGIN_DIR:-~/.local/share/tykctl/portal/plugins}"

echo "Installing $PLUGIN_NAME..."

# Create plugin directory
mkdir -p "$INSTALL_DIR"

# Copy plugin
cp "$PLUGIN_FILE" "$INSTALL_DIR/"
chmod +x "$INSTALL_DIR/$PLUGIN_FILE"

echo "Plugin installed successfully!"
echo "Usage: tykctl-portal $PLUGIN_NAME [options]"
```

## 🚀 Best Practices

### 1. Error Handling

```bash
# Good: Check required environment variables
if [[ -z "${TYK_PORTAL_URL:-}" ]]; then
    echo "Error: TYK_PORTAL_URL not set" >&2
    exit 1
fi

# Good: Use appropriate exit codes
if ! command -v curl >/dev/null 2>&1; then
    echo "Error: curl not found" >&2
    exit 2
fi
```

### 2. Logging

```bash
# Good: Consistent logging format
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] [$PLUGIN_NAME] $*" >&2
}

# Good: Different log levels
log_info() { log "INFO: $*"; }
log_warn() { log "WARN: $*"; }
log_error() { log "ERROR: $*"; }
```

### 3. Configuration

```bash
# Good: Support multiple configuration methods
CONFIG_FILE="${TYKCTL_PORTAL_CONFIG_FILE:-config.yaml}"
DRY_RUN="${TYKCTL_PORTAL_DRY_RUN:-false}"
TIMEOUT="${TYKCTL_PORTAL_TIMEOUT:-30s}"
```

### 4. Cross-Platform Compatibility

```bash
# Good: Detect platform
case "$(uname -s)" in
    Linux*)     PLATFORM="linux";;
    Darwin*)    PLATFORM="macos";;
    CYGWIN*|MINGW*|MSYS*) PLATFORM="windows";;
    *)          PLATFORM="unknown";;
esac
```

## 🆘 Troubleshooting

### Common Issues

**Plugin not found:**
- Check naming convention
- Verify executable permissions
- Check plugin discovery paths

**Permission denied:**
- Ensure plugin is executable
- Check directory permissions
- Verify user has access

**Environment variables not set:**
- Check extension configuration
- Verify environment setup
- Review plugin documentation

### Debug Mode

```bash
# Enable debug logging
export TYKCTL_DEBUG=true
export TYKCTL_VERBOSE=true

# Run plugin with debug info
tykctl-portal my-plugin --verbose
```

## 📚 Resources

- [Plugin System Documentation](../plugin/README.md)
- [Configuration Guide](../config/README.md)
- [API Documentation](../api/README.md)
- [Getting Started Guide](getting-started.md)