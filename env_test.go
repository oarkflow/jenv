package jenv_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

func TestUnmarshalJSON(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVICE_NAME", "MyTestService")
	os.Setenv("ENABLE_SERVICE", "false")
	os.Setenv("TIMEOUT", "15s")
	os.Setenv("START_TIME", "2024-02-01T15:00:00Z")
	os.Setenv("RATE", "2.5")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "1234")
	os.Setenv("DB_REPLICA_PORT", "5678")

	// Input JSON with placeholders
	jsonData := []byte(`
	{
	    "service": {
	        "name": "${SERVICE_NAME:MyService}",
	        "enabled": "${ENABLE_SERVICE:true}",
	        "timeout": "${TIMEOUT:30s}",
	        "start_time": "${START_TIME}",
	        "rate": "${RATE:1}"
	    },
	    "database": {
	        "hosts": ["${DB_HOST:localhost}"],
	        "ports": {"primary": "${DB_PORT:5432}", "replica": "${DB_REPLICA_PORT:5433}"}
	    }
	}`)

	var config Config
	err := jenv.UnmarshalJSON(jsonData, &config)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, "MyTestService", config.Service.Name)
	assert.False(t, config.Service.Enabled)
	assert.Equal(t, 2.5, config.Service.Rate)
	assert.Equal(t, 15*time.Second, config.Service.Timeout)

	expectedStartTime, _ := time.Parse(time.RFC3339, "2024-02-01T15:00:00Z")
	assert.Equal(t, expectedStartTime, config.Service.StartTime)

	assert.Equal(t, []string{"db.example.com"}, config.Database.Hosts)
	assert.Equal(t, map[string]int{"primary": 1234, "replica": 5678}, config.Database.Ports)
}

func TestUnmarshalYAML(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVICE_NAME", "YamlTestService")
	os.Setenv("ENABLE_SERVICE", "true")
	os.Setenv("TIMEOUT", "45s")
	os.Setenv("START_TIME", "2024-03-01T12:00:00Z")
	os.Setenv("RATE", "3.14")
	os.Setenv("DB_HOST", "yaml-db.example.com")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_REPLICA_PORT", "3307")

	// Input YAML with placeholders
	yamlData := []byte(`
service:
  name: "${SERVICE_NAME:DefaultService}"
  enabled: "${ENABLE_SERVICE:false}"
  timeout: "${TIMEOUT:60s}"
  start_time: "${START_TIME:2023-01-01T00:00:00Z}"
  rate: "${RATE:1.0}"
database:
  hosts:
    - "${DB_HOST:localhost}"
  ports:
    primary: "${DB_PORT:5432}"
    replica: "${DB_REPLICA_PORT:5433}"
`)

	var config Config
	err := jenv.UnmarshalYAML(yamlData, &config)
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, "YamlTestService", config.Service.Name)
	assert.True(t, config.Service.Enabled)
	assert.Equal(t, 3.14, config.Service.Rate)
	assert.Equal(t, 45*time.Second, config.Service.Timeout)

	expectedStartTime, _ := time.Parse(time.RFC3339, "2024-03-01T12:00:00Z")
	assert.Equal(t, expectedStartTime, config.Service.StartTime)

	assert.Equal(t, []string{"yaml-db.example.com"}, config.Database.Hosts)
	assert.Equal(t, map[string]int{"primary": 3306, "replica": 3307}, config.Database.Ports)
}
