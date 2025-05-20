# go-yaml-env

A Go package for loading YAML configuration files with environment variable support.

## Features

- Load YAML configuration files
- Support environment variable substitution in YAML
- Default values for environment variables
- Automatic .env file loading
- Generic type support for configuration structs

## Installation

```bash
go get github.com/ifuryst/go-yaml-env
```

## Usage

1. Create a configuration struct:

```go
type Config struct {
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
```

2. Create a YAML configuration file (e.g., `config.yaml`):

```yaml
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
```

3. Load the configuration:

```go
import "github.com/ifuryst/go-yaml-env"

func main() {
    cfg, path, err := yamlenv.LoadConfig[Config]("config.yaml")
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // Use the configuration
    fmt.Printf("Logger level: %s\n", cfg.Logger.Level)
}
```

## Environment Variable Syntax

The package supports the following syntax for environment variables in YAML:

- `${VAR}` - Required environment variable
- `${VAR:default}` - Environment variable with default value

## License

[MIT](LICENSE)
