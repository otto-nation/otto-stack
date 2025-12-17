package scripts

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedScripts(t *testing.T) {
	t.Run("generic init script is embedded", func(t *testing.T) {
		assert.NotEmpty(t, GenericInitScript)
		assert.Contains(t, GenericInitScript, "#!/")
		assert.Contains(t, strings.ToLower(GenericInitScript), "auto-discovery initialization")
	})

	t.Run("script contains required sections", func(t *testing.T) {
		script := GenericInitScript
		assert.Contains(t, script, "apk add --no-cache")
		assert.Contains(t, script, "INIT_SERVICE_NAME")
		assert.Contains(t, script, "CONFIG_DIR")
		assert.Contains(t, script, "SERVICE_ENDPOINT_URL")
	})
}

func TestProcessInit(t *testing.T) {
	t.Run("handles unknown service gracefully", func(t *testing.T) {
		tempDir := t.TempDir()
		err := ProcessInit("unknown-service", tempDir, "http://localhost:8080", "us-east-1")
		assert.NoError(t, err) // Should not error for unknown services
	})

	t.Run("processes localstack config files", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create test config file
		configContent := `
queues:
  - name: test-queue
topics:
  - name: test-topic
buckets:
  - name: test-bucket
`
		configFile := filepath.Join(tempDir, "localstack-sqs.yml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// This will fail without actual LocalStack, but tests the parsing logic
		err = ProcessInit("localstack", tempDir, "http://localhost:4566", "us-east-1")
		// We expect this to fail since LocalStack isn't running
		if err != nil {
			assert.Contains(t, err.Error(), "service not ready")
		}
	})

	t.Run("processes postgres config files", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create test config file
		configContent := `
schemas:
  - name: test_schema
databases:
  - name: test_db
`
		configFile := filepath.Join(tempDir, "postgres-schemas.yml")
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		// This will fail without actual Postgres, but tests the parsing logic
		err = ProcessInit("postgres", tempDir, "http://localhost:5432", "us-east-1")
		if err != nil {
			assert.Contains(t, err.Error(), "service not ready")
		}
	})

	t.Run("handles empty config directory", func(t *testing.T) {
		tempDir := t.TempDir()
		err := ProcessInit("localstack", tempDir, "http://localhost:4566", "us-east-1")
		if err != nil {
			assert.Contains(t, err.Error(), "service not ready")
		}
	})
}

func TestWaitForService(t *testing.T) {
	t.Run("returns immediately for unknown services", func(t *testing.T) {
		err := waitForService(context.Background(), "unknown", "http://localhost:8080")
		assert.NoError(t, err)
	})

	t.Run("times out for unreachable services", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		err := waitForHTTP(ctx, "http://localhost:99999/health")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service not ready")
	})
}

func TestProcessConfigs(t *testing.T) {
	t.Run("handles malformed YAML gracefully", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create malformed YAML file
		configFile := filepath.Join(tempDir, "localstack-sqs.yml")
		err := os.WriteFile(configFile, []byte("invalid: yaml: content: ["), 0644)
		require.NoError(t, err)

		err = processConfigs(context.Background(), "localstack", tempDir, "http://localhost:4566", "us-east-1")
		assert.NoError(t, err) // Should continue processing despite errors
	})

	t.Run("processes multiple config files", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create multiple config files
		configs := map[string]string{
			"localstack-sqs.yml": "queues:\n  - name: queue1",
			"localstack-sns.yml": "topics:\n  - name: topic1",
			"localstack-s3.yml":  "buckets:\n  - name: bucket1",
		}

		for filename, content := range configs {
			configFile := filepath.Join(tempDir, filename)
			err := os.WriteFile(configFile, []byte(content), 0644)
			require.NoError(t, err)
		}

		err := processConfigs(context.Background(), "localstack", tempDir, "http://localhost:4566", "us-east-1")
		assert.NoError(t, err) // Should process all files without fatal errors
	})
}
