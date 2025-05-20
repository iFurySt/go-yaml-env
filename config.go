package yamlenv

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from a YAML file with environment variable support
func LoadConfig[T any](filename string) (*T, string, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	cfgPath := getCfgPath(filename)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, cfgPath, err
	}

	// Resolve environment variables
	data = resolveEnv(data)

	var cfg T
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, cfgPath, err
	}

	return &cfg, cfgPath, nil
}

// resolveEnv replaces environment variable placeholders in YAML content
func resolveEnv(content []byte) []byte {
	regex := regexp.MustCompile(`\$\{(\w+)(?::([^}]*))?\}`)

	return regex.ReplaceAllFunc(content, func(match []byte) []byte {
		matches := regex.FindSubmatch(match)
		envKey := string(matches[1])
		var defaultValue string

		if len(matches) > 2 {
			defaultValue = string(matches[2])
		}

		if value, exists := os.LookupEnv(envKey); exists {
			return []byte(value)
		}
		return []byte(defaultValue)
	})
}

// getCfgPath returns the absolute path of the configuration file
func getCfgPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}

	// Try to find the config file in the current directory first
	if _, err := os.Stat(filename); err == nil {
		absPath, _ := filepath.Abs(filename)
		return absPath
	}

	// If not found, try to find it in the config directory
	configDir := "config"
	if _, err := os.Stat(configDir); err == nil {
		absPath, _ := filepath.Abs(filepath.Join(configDir, filename))
		return absPath
	}

	// If still not found, return the absolute path of the original filename
	absPath, _ := filepath.Abs(filename)
	return absPath
}
