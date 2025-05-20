package yamlenv

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Logger struct {
		Level      string `yaml:"level"`
		Format     string `yaml:"format"`
		Output     string `yaml:"output"`
		MaxSize    int    `yaml:"max_size"`
		MaxBackups int    `yaml:"max_backups"`
		MaxAge     int    `yaml:"max_age"`
		Compress   bool   `yaml:"compress"`
		Color      bool   `yaml:"color"`
		Stacktrace bool   `yaml:"stacktrace"`
	} `yaml:"logger"`
}

func TestLoadConfig(t *testing.T) {
	// Create a temporary test config file
	testConfig := `
logger:
  level: "${LOGGER_LEVEL:info}"
  format: "${LOGGER_FORMAT:console}"
  output: "${LOGGER_OUTPUT:stdout}"
  max_size: ${LOGGER_MAX_SIZE:100}
  max_backups: ${LOGGER_MAX_BACKUPS:3}
  max_age: ${LOGGER_MAX_AGE:7}
  compress: ${LOGGER_COMPRESS:true}
  color: ${LOGGER_COLOR:true}
  stacktrace: ${LOGGER_STACKTRACE:true}
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	assert.NoError(t, err)

	// Test with default values
	cfg, path, err := LoadConfig[TestConfig](configPath)
	assert.NoError(t, err)
	assert.Equal(t, configPath, path)
	assert.Equal(t, "info", cfg.Logger.Level)
	assert.Equal(t, "console", cfg.Logger.Format)
	assert.Equal(t, "stdout", cfg.Logger.Output)
	assert.Equal(t, 100, cfg.Logger.MaxSize)
	assert.Equal(t, 3, cfg.Logger.MaxBackups)
	assert.Equal(t, 7, cfg.Logger.MaxAge)
	assert.True(t, cfg.Logger.Compress)
	assert.True(t, cfg.Logger.Color)
	assert.True(t, cfg.Logger.Stacktrace)

	// Test with environment variables
	os.Setenv("LOGGER_LEVEL", "debug")
	os.Setenv("LOGGER_FORMAT", "json")
	os.Setenv("LOGGER_MAX_SIZE", "200")
	os.Setenv("LOGGER_COMPRESS", "false")

	cfg, path, err = LoadConfig[TestConfig](configPath)
	assert.NoError(t, err)
	assert.Equal(t, configPath, path)
	assert.Equal(t, "debug", cfg.Logger.Level)
	assert.Equal(t, "json", cfg.Logger.Format)
	assert.Equal(t, "stdout", cfg.Logger.Output)
	assert.Equal(t, 200, cfg.Logger.MaxSize)
	assert.Equal(t, 3, cfg.Logger.MaxBackups)
	assert.Equal(t, 7, cfg.Logger.MaxAge)
	assert.False(t, cfg.Logger.Compress)
	assert.True(t, cfg.Logger.Color)
	assert.True(t, cfg.Logger.Stacktrace)

	// Clean up environment variables
	os.Unsetenv("LOGGER_LEVEL")
	os.Unsetenv("LOGGER_FORMAT")
	os.Unsetenv("LOGGER_MAX_SIZE")
	os.Unsetenv("LOGGER_COMPRESS")
}

func TestLoadConfigWithInvalidFile(t *testing.T) {
	absPath, _ := filepath.Abs("nonexistent.yaml")
	_, path, err := LoadConfig[TestConfig]("nonexistent.yaml")
	assert.Error(t, err)
	assert.Equal(t, absPath, path)
}

func TestLoadConfigWithInvalidYAML(t *testing.T) {
	// Create a temporary test config file with invalid YAML
	testConfig := `
logger:
  level: "${LOGGER_LEVEL:info}"
  format: "${LOGGER_FORMAT:console}"
  output: "${LOGGER_OUTPUT:stdout}"
  max_size: ${LOGGER_MAX_SIZE:100}
  max_backups: ${LOGGER_MAX_BACKUPS:3}
  max_age: ${LOGGER_MAX_AGE:7}
  compress: ${LOGGER_COMPRESS:true}
  color: ${LOGGER_COLOR:true}
  stacktrace: ${LOGGER_STACKTRACE:true}
invalid: yaml: here
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	assert.NoError(t, err)

	_, path, err := LoadConfig[TestConfig](configPath)
	assert.Error(t, err)
	assert.Equal(t, configPath, path)
}

func TestGetCfgPath(t *testing.T) {
	// Test absolute path
	absPath := "/absolute/path/config.yaml"
	assert.Equal(t, absPath, getCfgPath(absPath))

	// Test relative path
	relPath := "config.yaml"
	expectedPath, _ := filepath.Abs(relPath)
	assert.Equal(t, expectedPath, getCfgPath(relPath))
}

func TestLoadConfigWithDotEnv(t *testing.T) {
	// 清理相关环境变量，确保只从.env读取
	os.Unsetenv("LOGGER_LEVEL")
	os.Unsetenv("LOGGER_FORMAT")

	tmpDir := t.TempDir()

	dotenvPath := filepath.Join(tmpDir, ".env")
	dotenvContent := "LOGGER_LEVEL=warn\nLOGGER_FORMAT=logfmt\n"
	err := os.WriteFile(dotenvPath, []byte(dotenvContent), 0644)
	assert.NoError(t, err)

	testConfig := `
logger:
  level: "${LOGGER_LEVEL:info}"
  format: "${LOGGER_FORMAT:console}"
`
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	err = os.WriteFile(configPath, []byte(testConfig), 0644)
	assert.NoError(t, err)

	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	cfg, path, err := LoadConfig[TestConfig](configPath)
	assert.NoError(t, err)
	assert.Equal(t, configPath, path)
	assert.Equal(t, "warn", cfg.Logger.Level)
	assert.Equal(t, "logfmt", cfg.Logger.Format)
}
