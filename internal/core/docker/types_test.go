//go:build unit

package docker

import (
	"testing"

	"github.com/otto-nation/otto-stack/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDockerConstants(t *testing.T) {
	t.Run("validates docker constants are defined", func(t *testing.T) {
		constants := []string{
			ComposeFieldServices,
			ComposeFieldNetworks,
			ComposeFieldVolumes,
		}

		for _, constant := range constants {
			assert.NotEmpty(t, constant, "Docker constant should not be empty")
		}
	})

	t.Run("validates compose file path uses core constants", func(t *testing.T) {
		assert.Contains(t, DockerComposeFilePath, core.OttoStackDir)
		assert.Contains(t, DockerComposeFilePath, DockerComposeFileName)
	})
}

func TestInitServiceSpec_Validation(t *testing.T) {
	t.Run("validates init service spec structure", func(t *testing.T) {
		spec := InitServiceSpec{
			Image:   "postgres:13",
			Enabled: true,
			Environment: map[string]string{
				"POSTGRES_DB": "testdb",
			},
		}

		assert.Equal(t, "postgres:13", spec.Image)
		assert.True(t, spec.Enabled)
		assert.Equal(t, "testdb", spec.Environment["POSTGRES_DB"])
	})

	t.Run("handles empty init service spec", func(t *testing.T) {
		spec := InitServiceSpec{}

		assert.Empty(t, spec.Image)
		assert.False(t, spec.Enabled)
		assert.Empty(t, spec.Environment)
	})
}

func TestInitScript_Validation(t *testing.T) {
	t.Run("validates init script structure", func(t *testing.T) {
		script := InitScript{
			Content: "CREATE DATABASE test;",
		}

		assert.Equal(t, "CREATE DATABASE test;", script.Content)
	})

	t.Run("validates init script type constants", func(t *testing.T) {
		types := []string{
			InitScriptTypeShell,
			InitScriptTypeSQL,
			InitScriptTypeAWSResources,
			InitScriptTypeKafkaTopics,
		}

		for _, scriptType := range types {
			assert.NotEmpty(t, scriptType, "Init script type should not be empty")
		}
	})
}

func TestResourceType_Constants(t *testing.T) {
	t.Run("validates resource type constants", func(t *testing.T) {
		types := []ResourceType{
			ResourceContainer,
			ResourceNetwork,
			ResourceVolume,
		}

		for _, resourceType := range types {
			assert.NotEmpty(t, string(resourceType), "Resource type should not be empty")
		}
	})
}
