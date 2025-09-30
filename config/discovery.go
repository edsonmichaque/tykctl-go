package config

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adrg/xdg"
)

// DiscoverHooks discovers hooks with enhanced filtering and metadata
func DiscoverHooks(extension string, filter HookFilter, contextPaths []string) ([]Hook, error) {
	var hooks []Hook

	paths := getHookSearchPaths(extension, contextPaths)

	for _, path := range paths {
		if entries, err := os.ReadDir(path); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && entry.Type()&0111 != 0 {
					hook := Hook{
						Name:         entry.Name(),
						Path:         filepath.Join(path, entry.Name()),
						Event:        extractEventFromName(entry.Name()),
						Priority:     extractPriorityFromName(entry.Name()),
						Timeout:      extractTimeoutFromName(entry.Name()),
						Enabled:      true,
						Metadata:     make(map[string]string),
						DiscoveredAt: time.Now(),
					}

					// Calculate checksum
					if checksum, err := calculateFileChecksum(hook.Path); err == nil {
						hook.Checksum = checksum
					}

					// Apply filters
					if matchesHookFilter(hook, filter) {
						hooks = append(hooks, hook)
					}
				}
			}
		}
	}

	return hooks, nil
}

// DiscoverPlugins discovers plugins with enhanced filtering and metadata
func DiscoverPlugins(extension string, filter PluginFilter, contextPaths []string) ([]Plugin, error) {
	var plugins []Plugin

	paths := getPluginSearchPaths(extension, contextPaths)
	prefix := fmt.Sprintf("tykctl-%s-", extension)

	for _, path := range paths {
		if entries, err := os.ReadDir(path); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && entry.Type()&0111 != 0 {
					if strings.HasPrefix(entry.Name(), prefix) {
						plugin := Plugin{
							Name:         strings.TrimPrefix(entry.Name(), prefix),
							Path:         filepath.Join(path, entry.Name()),
							Version:      extractVersionFromPath(entry.Name()),
							Commands:     extractCommandsFromPath(entry.Name()),
							Enabled:      true,
							Metadata:     make(map[string]string),
							DiscoveredAt: time.Now(),
						}

						// Calculate checksum
						if checksum, err := calculateFileChecksum(plugin.Path); err == nil {
							plugin.Checksum = checksum
						}

						// Apply filters
						if matchesPluginFilter(plugin, filter) {
							plugins = append(plugins, plugin)
						}
					}
				}
			}
		}
	}

	return plugins, nil
}

// DiscoverTemplates discovers templates with enhanced filtering and metadata
func DiscoverTemplates(extension string, filter TemplateFilter, contextPaths []string) ([]Template, error) {
	var templates []Template

	paths := getTemplateSearchPaths(extension, contextPaths)

	for _, path := range paths {
		if entries, err := os.ReadDir(path); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					template := Template{
						Name:         entry.Name(),
						Path:         filepath.Join(path, entry.Name()),
						Type:         extractTypeFromName(entry.Name()),
						Format:       extractFormatFromName(entry.Name()),
						Schema:       extractSchemaFromPath(entry.Name()),
						Enabled:      true,
						Metadata:     make(map[string]string),
						DiscoveredAt: time.Now(),
					}

					// Calculate checksum
					if checksum, err := calculateFileChecksum(template.Path); err == nil {
						template.Checksum = checksum
					}

					// Apply filters
					if matchesTemplateFilter(template, filter) {
						templates = append(templates, template)
					}
				}
			}
		}
	}

	return templates, nil
}

// DiscoverCache discovers cache configurations for an extension
func DiscoverCache(extension string, contextPaths []string) ([]CacheConfig, error) {
	var cacheConfigs []CacheConfig

	paths := getCacheSearchPaths(extension, contextPaths)

	for _, path := range paths {
		if entries, err := os.ReadDir(path); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					cacheConfig := CacheConfig{
						Name:         entry.Name(),
						Path:         filepath.Join(path, entry.Name()),
						Type:         extractCacheTypeFromName(entry.Name()),
						Enabled:      true,
						Metadata:     make(map[string]string),
						DiscoveredAt: time.Now(),
					}

					// Calculate checksum
					if checksum, err := calculateFileChecksum(cacheConfig.Path); err == nil {
						cacheConfig.Checksum = checksum
					}

					cacheConfigs = append(cacheConfigs, cacheConfig)
				}
			}
		}
	}

	return cacheConfigs, nil
}

// Helper functions
func calculateFileChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

func extractEventFromName(name string) string {
	// Extract event from filename like "pre-command-setup.sh"
	parts := strings.Split(name, "-")
	if len(parts) >= 2 {
		return parts[0] + "-" + parts[1]
	}
	return "unknown"
}

func extractPriorityFromName(name string) int {
	// Extract priority from filename like "100-pre-command-setup.sh"
	parts := strings.Split(name, "-")
	if len(parts) > 0 {
		if priority, err := strconv.Atoi(parts[0]); err == nil {
			return priority
		}
	}
	return 100 // Default priority
}

func extractTimeoutFromName(name string) time.Duration {
	// Extract timeout from filename like "30s-pre-command-setup.sh"
	parts := strings.Split(name, "-")
	if len(parts) > 0 {
		if timeout, err := time.ParseDuration(parts[0]); err == nil {
			return timeout
		}
	}
	return 30 * time.Second // Default timeout
}

func extractVersionFromPath(name string) string {
	// Extract version from filename like "tykctl-ai-studio-trainer-v1.0.0"
	parts := strings.Split(name, "-")
	for i, part := range parts {
		if strings.HasPrefix(part, "v") && len(part) > 1 {
			return part
		}
		if i > 0 && strings.Contains(part, ".") {
			return part
		}
	}
	return "unknown"
}

func extractCommandsFromPath(name string) []string {
	// Extract commands from filename like "tykctl-ai-studio-trainer"
	parts := strings.Split(name, "-")
	if len(parts) >= 3 {
		return []string{parts[len(parts)-1]}
	}
	return []string{}
}

func extractTypeFromName(name string) string {
	// Extract type from filename like "output-model.json"
	parts := strings.Split(name, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func extractFormatFromName(name string) string {
	// Extract format from filename like "output-model.json"
	ext := filepath.Ext(name)
	if ext != "" {
		return ext[1:] // Remove the dot
	}
	return "unknown"
}

func extractSchemaFromPath(name string) string {
	// Extract schema from filename like "output-model-schema.json"
	if strings.Contains(name, "schema") {
		return "json-schema"
	}
	return ""
}

func extractCacheTypeFromName(name string) string {
	// Extract cache type from filename like "redis-cache.yaml"
	parts := strings.Split(name, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func getHookSearchPaths(extension string, contextPaths []string) []string {
	paths := []string{
		filepath.Join("/etc/tykctl", extension, "hooks"),
		filepath.Join(xdg.DataHome, "tykctl", extension, "hooks"),
		filepath.Join(os.Getenv("HOME"), ".tykctl", extension, "hooks"),
		filepath.Join(".tykctl", extension, "hooks"),
	}
	
	// Add context-specific paths
	for _, contextPath := range contextPaths {
		paths = append(paths, filepath.Join(contextPath, extension, "hooks"))
	}
	
	return paths
}

func getPluginSearchPaths(extension string, contextPaths []string) []string {
	paths := []string{
		"/usr/local/bin",
		"/usr/bin",
		filepath.Join(xdg.DataHome, "tykctl", extension, "bin"),
		filepath.Join(os.Getenv("HOME"), ".tykctl", extension, "bin"),
		filepath.Join(".tykctl", extension, "bin"),
	}
	
	// Add context-specific paths
	for _, contextPath := range contextPaths {
		paths = append(paths, filepath.Join(contextPath, extension, "bin"))
	}
	
	return paths
}

func getTemplateSearchPaths(extension string, contextPaths []string) []string {
	paths := []string{
		filepath.Join("/etc/tykctl", extension, "templates"),
		filepath.Join(xdg.DataHome, "tykctl", extension, "templates"),
		filepath.Join(os.Getenv("HOME"), ".tykctl", extension, "templates"),
		filepath.Join(".tykctl", extension, "templates"),
	}
	
	// Add context-specific paths
	for _, contextPath := range contextPaths {
		paths = append(paths, filepath.Join(contextPath, extension, "templates"))
	}
	
	return paths
}

func getCacheSearchPaths(extension string, contextPaths []string) []string {
	paths := []string{
		filepath.Join("/etc/tykctl", extension, "cache"),
		filepath.Join(xdg.DataHome, "tykctl", extension, "cache"),
		filepath.Join(os.Getenv("HOME"), ".tykctl", extension, "cache"),
		filepath.Join(".tykctl", extension, "cache"),
	}
	
	// Add context-specific paths
	for _, contextPath := range contextPaths {
		paths = append(paths, filepath.Join(contextPath, extension, "cache"))
	}
	
	return paths
}

func matchesHookFilter(hook Hook, filter HookFilter) bool {
	// Apply event filter
	if len(filter.Events) > 0 {
		found := false
		for _, event := range filter.Events {
			if hook.Event == event {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Apply enabled filter
	if filter.Enabled != nil && hook.Enabled != *filter.Enabled {
		return false
	}

	// Apply priority filter
	if filter.Priority != nil && hook.Priority != *filter.Priority {
		return false
	}

	// Apply timeout filters
	if filter.MinTimeout > 0 && hook.Timeout < filter.MinTimeout {
		return false
	}
	if filter.MaxTimeout > 0 && hook.Timeout > filter.MaxTimeout {
		return false
	}

	// Apply metadata filter
	for key, value := range filter.Metadata {
		if hook.Metadata[key] != value {
			return false
		}
	}

	return true
}

func matchesPluginFilter(plugin Plugin, filter PluginFilter) bool {
	// Apply commands filter
	if len(filter.Commands) > 0 {
		found := false
		for _, command := range filter.Commands {
			for _, pluginCommand := range plugin.Commands {
				if pluginCommand == command {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// Apply enabled filter
	if filter.Enabled != nil && plugin.Enabled != *filter.Enabled {
		return false
	}

	// Apply version filter
	if filter.Version != "" && plugin.Version != filter.Version {
		return false
	}

	// Apply metadata filter
	for key, value := range filter.Metadata {
		if plugin.Metadata[key] != value {
			return false
		}
	}

	return true
}

func matchesTemplateFilter(template Template, filter TemplateFilter) bool {
	// Apply types filter
	if len(filter.Types) > 0 {
		found := false
		for _, templateType := range filter.Types {
			if template.Type == templateType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Apply formats filter
	if len(filter.Formats) > 0 {
		found := false
		for _, format := range filter.Formats {
			if template.Format == format {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Apply enabled filter
	if filter.Enabled != nil && template.Enabled != *filter.Enabled {
		return false
	}

	// Apply metadata filter
	for key, value := range filter.Metadata {
		if template.Metadata[key] != value {
			return false
		}
	}

	return true
}