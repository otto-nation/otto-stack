package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testServiceName = "test-service"
	testEndpointURL = "http://test-service:1234"
	testConfigDir   = "/test-config"
	defaultHost     = core.ServiceLocalhost
	defaultPort     = core.PortPostgreSQL
	defaultDatabase = "testdb"
	customHost      = "custom-host"
	customPort      = "9999"
	customDatabase  = "customdb"
)

// mockServiceConfig implements ServiceConfigInterface for testing
type mockServiceConfig struct {
	initImage      string
	environment    map[string]string
	connectionPort int
	hasConnection  bool
}

func (m *mockServiceConfig) GetInitContainerImage() string {
	return m.initImage
}

func (m *mockServiceConfig) GetEnvironment() map[string]string {
	return m.environment
}

func (m *mockServiceConfig) GetConnectionPort() int {
	return m.connectionPort
}

func (m *mockServiceConfig) HasConnection() bool {
	return m.hasConnection
}

func (m *mockServiceConfig) GetInitContainerSpec() *InitContainerSpec {
	return nil // Return nil for tests that don't need init container spec
}

func TestInitContainerManager_TemplateResolution(t *testing.T) {
	manager := &InitContainerManager{}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple_default_resolution",
			input:    "${HOST:-" + defaultHost + "}",
			expected: defaultHost,
		},
		{
			name:     "port_default_resolution",
			input:    "${PORT:-" + defaultPort + "}",
			expected: defaultPort,
		},
		{
			name:     "database_default_resolution",
			input:    "${DATABASE:-" + defaultDatabase + "}",
			expected: defaultDatabase,
		},
		{
			name:     "complex_url_resolution",
			input:    "postgresql://${USER:-postgres}:${PASSWORD:-password}@${HOST:-" + defaultHost + "}:${PORT:-" + defaultPort + "}/${DB:-" + defaultDatabase + "}",
			expected: "postgresql://postgres:password@" + defaultHost + ":" + defaultPort + "/" + defaultDatabase,
		},
		{
			name:     "no_template_unchanged",
			input:    "plain-string-value",
			expected: "plain-string-value",
		},
		{
			name:     "empty_string_unchanged",
			input:    "",
			expected: "",
		},
		{
			name:     "multiple_templates_in_string",
			input:    "http://${HOST:-" + defaultHost + "}:${PORT:-" + defaultPort + "}/path",
			expected: "http://" + defaultHost + ":" + defaultPort + "/path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := manager.resolveTemplate(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInitContainerManager_ApplyServiceEnvironment(t *testing.T) {
	manager := &InitContainerManager{}

	t.Run("copies_and_resolves_service_environment", func(t *testing.T) {
		service := &mockServiceConfig{
			environment: map[string]string{
				"HOST":     "${HOST:-" + defaultHost + "}",
				"PORT":     "${PORT:-" + defaultPort + "}",
				"DATABASE": "${DATABASE:-" + defaultDatabase + "}",
				"PLAIN":    "plain-value",
			},
		}

		config := &InitContainerConfig{
			Environment: map[string]string{
				InitServiceName: testServiceName,
			},
		}

		manager.applyServiceEnvironment(config, service)

		// Verify template resolution
		assert.Equal(t, defaultHost, config.Environment["HOST"])
		assert.Equal(t, defaultPort, config.Environment["PORT"])
		assert.Equal(t, defaultDatabase, config.Environment["DATABASE"])
		assert.Equal(t, "plain-value", config.Environment["PLAIN"])

		// Verify original environment preserved
		assert.Equal(t, testServiceName, config.Environment[InitServiceName])
	})

	t.Run("sets_endpoint_url_when_connection_available", func(t *testing.T) {
		service := &mockServiceConfig{
			environment:    map[string]string{},
			hasConnection:  true,
			connectionPort: 5432,
		}

		config := &InitContainerConfig{
			Environment: map[string]string{
				InitServiceName: testServiceName,
			},
		}

		manager.applyServiceEnvironment(config, service)

		expectedURL := "http://" + testServiceName + ":" + core.PortPostgreSQL
		assert.Equal(t, expectedURL, config.Environment[InitServiceEndpointURL])
	})
}

func TestInitContainerManager_GetInitContainerImage(t *testing.T) {
	manager := &InitContainerManager{}

	testCases := []struct {
		name     string
		service  ServiceConfigInterface
		expected string
	}{
		{
			name: "uses_custom_init_image",
			service: &mockServiceConfig{
				initImage: "custom/init:latest",
			},
			expected: "custom/init:latest",
		},
		{
			name: "uses_default_when_no_init_image",
			service: &mockServiceConfig{
				initImage: "",
			},
			expected: AlpineLatestImage,
		},
		{
			name: "uses_default_when_empty_init_image",
			service: &mockServiceConfig{
				initImage: "",
			},
			expected: AlpineLatestImage,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := manager.getInitContainerImage(tc.service)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestInitContainerManager_BuildBaseInitConfig(t *testing.T) {
	manager := &InitContainerManager{}

	t.Run("creates_correct_base_config", func(t *testing.T) {
		service := &mockServiceConfig{
			initImage: "test/init:v1",
		}

		config := manager.buildBaseInitConfig(testServiceName, service, "test-project")

		require.NotNil(t, config)
		assert.Equal(t, "test/init:v1", config.Image)
		assert.Equal(t, testServiceName, config.Environment[InitServiceName])
		assert.Equal(t, "/config", config.Environment[InitConfigDir])
		assert.Contains(t, config.Networks[0], "test-project")
		assert.Contains(t, config.Networks[0], NetworkNameSuffix)
	})

	t.Run("uses_default_image_when_not_specified", func(t *testing.T) {
		service := &mockServiceConfig{
			initImage: "",
		}

		config := manager.buildBaseInitConfig(testServiceName, service, "test-project")

		assert.Equal(t, AlpineLatestImage, config.Image)
	})
}
