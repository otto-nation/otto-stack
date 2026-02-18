package fixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
	"gopkg.in/yaml.v3"
)

func getFixturePath(category, name string) string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return filepath.Join(dir, category, name+".yaml")
}

// LoadService loads a service config from test fixtures
func LoadService(t *testing.T, name string) servicetypes.ServiceConfig {
	t.Helper()

	data, err := os.ReadFile(getFixturePath("services", name))
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}

	var cfg servicetypes.ServiceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to parse fixture %s: %v", name, err)
	}

	return cfg
}

// LoadConfigYAML loads raw config YAML from test fixtures
func LoadConfigYAML(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(getFixturePath("configs", name))
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}

	return data
}
