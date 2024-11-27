# jenv

`jenv` is a Go package that simplifies configuration parsing by allowing placeholders in JSON and YAML files to be resolved dynamically using environment variables. The package supports default values and type-safe conversion of fields.

## Features
* Parse JSON and YAML configurations with environment variable resolution.
* Support for default values in ${VAR:default} syntax.
* Type-safe mapping of configuration values to Go structs.
* Handle complex data types such as time.Time, time.Duration, slices, and maps.

## Installation

> go get github.com/oarkflow/jenv

## Usage
### Environment Variables

Set environment variables for your configuration:

```bash
export SERVICE_NAME="MyAppService"
export ENABLE_SERVICE="true"
export RATE="2.5"
export TIMEOUT="15s"
export START_TIME="2024-01-01T12:00:00Z"
export DB_HOST="db.example.com"
export DB_PORT="5432"
export DB_REPLICA_PORT="5433"

```

## Example JSON Configuration
Create a JSON configuration file config.json:

```json
{
    "service": {
        "name": "${SERVICE_NAME:DefaultService}",
        "enabled": "${ENABLE_SERVICE:false}",
        "timeout": "${TIMEOUT:30s}",
        "start_time": "${START_TIME:2023-01-01T00:00:00Z}",
        "rate": "${RATE:1.0}"
    },
    "database": {
        "hosts": ["${DB_HOST:localhost}"],
        "ports": {
            "primary": "${DB_PORT:3306}",
            "replica": "${DB_REPLICA_PORT:3307}"
        }
    }
}

```

Parse it in your Go code:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/oarkflow/jenv"
)

type Service struct {
	Name      string        `json:"name"`
	Enabled   bool          `json:"enabled"`
	Rate      float64       `json:"rate"`
	Timeout   time.Duration `json:"timeout"`
	StartTime time.Time     `json:"start_time"`
}

type Database struct {
	Hosts []string       `json:"hosts"`
	Ports map[string]int `json:"ports"`
}

type Config struct {
	Service  Service  `json:"service"`
	Database Database `json:"database"`
}

func main() {
	os.Setenv("SERVICE_NAME", "MyAppService")
	os.Setenv("ENABLE_SERVICE", "true")
	os.Setenv("RATE", "2.5")
	os.Setenv("TIMEOUT", "15s")
	os.Setenv("START_TIME", "2024-01-01T12:00:00Z")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_REPLICA_PORT", "5433")

	jsonData := []byte(`
	{
	    "service": {
	        "name": "${SERVICE_NAME:DefaultService}",
	        "enabled": "${ENABLE_SERVICE:false}",
	        "timeout": "${TIMEOUT:30s}",
	        "start_time": "${START_TIME:2023-01-01T00:00:00Z}",
	        "rate": "${RATE:1.0}"
	    },
	    "database": {
	        "hosts": ["${DB_HOST:localhost}"],
	        "ports": {"primary": "${DB_PORT:3306}", "replica": "${DB_REPLICA_PORT:3307}"}
	    }
	}`)

	var config Config
	err := jenv.UnmarshalJSON(jsonData, &config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Config: %+v\n", config)
}

```

## Example YAML Configuration
Create a YAML configuration file config.yaml:

```yaml
service:
  name: "${SERVICE_NAME:DefaultService}"
  enabled: "${ENABLE_SERVICE:false}"
  timeout: "${TIMEOUT:30s}"
  start_time: "${START_TIME:2023-01-01T00:00:00Z}"
  rate: "${RATE:1.0}"
database:
  hosts:
    - "${DB_HOST:localhost}"
  ports:
    primary: "${DB_PORT:3306}"
    replica: "${DB_REPLICA_PORT:3307}"

```

Parse it in your Go code:
```go
yamlData := []byte(`
service:
  name: "${SERVICE_NAME:DefaultService}"
  enabled: "${ENABLE_SERVICE:false}"
  timeout: "${TIMEOUT:30s}"
  start_time: "${START_TIME:2023-01-01T00:00:00Z}"
  rate: "${RATE:1.0}"
database:
  hosts:
    - "${DB_HOST:localhost}"
  ports:
    primary: "${DB_PORT:3306}"
    replica: "${DB_REPLICA_PORT:3307}"
`)

var config Config
err := jenv.UnmarshalYAML(yamlData, &config)
if err != nil {
	fmt.Printf("Error: %v\n", err)
	return
}

fmt.Printf("Config: %+v\n", config)

```

Ensure you have environment variables set for the tests, or mock them in your test code.

## Contributing
We welcome contributions! Please follow these steps:

* Fork the repository.
* Create a feature branch: `git checkout -b feature-name`.
* Commit your changes: `git commit -m "Add feature"`.
* Push to the branch: `git push origin feature-name`.
* Open a Pull Request.
