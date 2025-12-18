package stack

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/otto-nation/otto-stack/internal/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpHandler_InitContainer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("creates valid init container config", func(t *testing.T) {
		// Setup test environment
		tempDir := t.TempDir()

		// Create mock config directory
		configDir := filepath.Join(tempDir, core.OttoStackDir, core.ServiceConfigsDir)
		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		// Create test config file
		configContent := `
queues:
  - name: test-queue
    attributes:
      VisibilityTimeout: "30"
topics:
  - name: test-topic
buckets:
  - name: test-bucket
`
		configFile := filepath.Join(configDir, "localstack-sqs.yml")
		err = os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// Change to temp directory
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		// Create handler and setup

		// Mock config
		cfg := &config.Config{
			Project: config.ProjectConfig{
				Name: "test-project",
			},
		}

		setup := &CoreSetup{
			Config: cfg,
		}

		// Test init container config creation
		initManager := NewInitContainerManager()
		initConfig, err := initManager.CreateInitContainerConfig("localstack", setup)
		assert.NoError(t, err)

		// Verify config
		assert.Equal(t, "localstack/localstack:latest", initConfig.Image)
		assert.Contains(t, initConfig.Command, "sh")
		assert.Contains(t, initConfig.Command, "-c")

		// Verify environment variables
		assert.Equal(t, "test", initConfig.Environment["AWS_ACCESS_KEY_ID"])
		assert.Equal(t, "us-east-1", initConfig.Environment["AWS_DEFAULT_REGION"])
		assert.Equal(t, "localstack", initConfig.Environment["INIT_SERVICE_NAME"])
		assert.Equal(t, "/config", initConfig.Environment["CONFIG_DIR"])
		assert.Contains(t, initConfig.Environment["SERVICE_ENDPOINT_URL"], "localhost:4566")

		// Verify volumes
		assert.Len(t, initConfig.Volumes, 1)
		assert.Contains(t, initConfig.Volumes[0], "/config")

		// Verify networks
		assert.Len(t, initConfig.Networks, 1)
		assert.Equal(t, "test-project-network", initConfig.Networks[0])
	})

	t.Run("init script contains required tools installation", func(t *testing.T) {

		cfg := &config.Config{
			Project: config.ProjectConfig{Name: "test"},
		}
		setup := &CoreSetup{Config: cfg}

		initConfig, err := NewInitContainerManager().createInitContainerConfig("localstack", setup)
		assert.NoError(t, err)

		// The script should be in the command
		scriptContent := initConfig.Command[2] // sh -c "script"

		assert.Contains(t, scriptContent, "apk add --no-cache")
		assert.Contains(t, scriptContent, "aws-cli")
		assert.Contains(t, scriptContent, "postgresql-client")
		assert.Contains(t, scriptContent, "INIT_SERVICE_NAME")
	})

	t.Run("handles different service types", func(t *testing.T) {
		cfg := &config.Config{
			Project: config.ProjectConfig{Name: "test"},
		}
		setup := &CoreSetup{Config: cfg}

		services := []string{"localstack", "postgres", "kafka"}

		for _, service := range services {
			t.Run(service, func(t *testing.T) {
				initConfig, err := NewInitContainerManager().createInitContainerConfig(service, setup)
				assert.NoError(t, err)

				assert.Equal(t, service, initConfig.Environment["INIT_SERVICE_NAME"])

				// Check service-specific endpoint URLs
				switch service {
				case "localstack":
					assert.Contains(t, initConfig.Environment["SERVICE_ENDPOINT_URL"], "localhost:4566")
				case "postgres":
					assert.Contains(t, initConfig.Environment["SERVICE_ENDPOINT_URL"], "postgresql://")
				case "kafka":
					assert.Contains(t, initConfig.Environment["SERVICE_ENDPOINT_URL"], service)
				}

				// Verify service-specific setup in script
				scriptContent := initConfig.Command[2]
				if service == "kafka" {
					assert.Contains(t, scriptContent, "kafka")
				}
			})
		}
	})
}

func TestUpHandler_InitContainer_ConfigDiscovery(t *testing.T) {
	t.Run("discovers multiple config files", func(t *testing.T) {
		tempDir := t.TempDir()
		configDir := filepath.Join(tempDir, core.OttoStackDir, core.ServiceConfigsDir)
		err := os.MkdirAll(configDir, 0755)
		require.NoError(t, err)

		// Create multiple config files
		configs := map[string]string{
			"localstack-sqs.yml": "queues:\n  - name: queue1",
			"localstack-sns.yml": "topics:\n  - name: topic1",
			"localstack-s3.yml":  "buckets:\n  - name: bucket1",
		}

		for filename, content := range configs {
			configFile := filepath.Join(configDir, filename)
			err := os.WriteFile(configFile, []byte(content), 0644)
			require.NoError(t, err)
		}

		// Change to temp directory
		originalDir, _ := os.Getwd()
		defer os.Chdir(originalDir)
		os.Chdir(tempDir)

		cfg := &config.Config{
			Project: config.ProjectConfig{Name: "test"},
		}
		setup := &CoreSetup{Config: cfg}

		initConfig, err := NewInitContainerManager().createInitContainerConfig("localstack", setup)
		assert.NoError(t, err)

		// Verify the config directory is mounted
		assert.Len(t, initConfig.Volumes, 1)
		assert.Contains(t, initConfig.Volumes[0], "/config")
	})
}
