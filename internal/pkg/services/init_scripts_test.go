//go:build unit

package services

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	servicetypes "github.com/otto-nation/otto-stack/internal/pkg/types"
)

func TestHasInitScripts(t *testing.T) {
	s := &Service{logger: logger.GetLogger()}

	t.Run("returns true when init service enabled", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test",
			InitService: &docker.InitServiceSpec{
				Enabled: true,
			},
		}
		assert.True(t, s.hasInitScripts(config))
	})

	t.Run("returns false when init service disabled", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test",
			InitService: &docker.InitServiceSpec{
				Enabled: false,
			},
		}
		assert.False(t, s.hasInitScripts(config))
	})

	t.Run("returns false when init service nil", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test",
		}
		assert.False(t, s.hasInitScripts(config))
	})
}

func TestHasLocalInitScripts(t *testing.T) {
	s := &Service{logger: logger.GetLogger()}

	t.Run("returns true for local mode", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test",
			InitService: &docker.InitServiceSpec{
				Enabled: true,
				Mode:    docker.InitServiceModeLocal,
			},
		}
		assert.True(t, s.hasLocalInitScripts(config))
	})

	t.Run("returns false for container mode", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test",
			InitService: &docker.InitServiceSpec{
				Enabled: true,
				Mode:    docker.InitServiceModeContainer,
			},
		}
		assert.False(t, s.hasLocalInitScripts(config))
	})

	t.Run("returns false when disabled", func(t *testing.T) {
		config := servicetypes.ServiceConfig{
			Name: "test",
			InitService: &docker.InitServiceSpec{
				Enabled: false,
				Mode:    docker.InitServiceModeLocal,
			},
		}
		assert.False(t, s.hasLocalInitScripts(config))
	})
}

func TestIsConfigStructField(t *testing.T) {
	type TestConfig struct {
		Value string
	}

	type TestStruct struct {
		Config *TestConfig
		Name   string
	}

	t.Run("returns true for pointer to struct with Config in name", func(t *testing.T) {
		s := TestStruct{}
		v := reflect.ValueOf(&s).Elem()
		field := v.Field(0)
		fieldType := v.Type().Field(0)
		assert.True(t, isConfigStructField(field, fieldType))
	})

	t.Run("returns false for non-pointer field", func(t *testing.T) {
		s := TestStruct{}
		v := reflect.ValueOf(&s).Elem()
		field := v.Field(1)
		fieldType := v.Type().Field(1)
		assert.False(t, isConfigStructField(field, fieldType))
	})
}

func TestGetYAMLFieldNameFromTag(t *testing.T) {
	type TestStruct struct {
		Field1 string `yaml:"field_one"`
		Field2 string `yaml:"field_two,omitempty"`
		Field3 string `yaml:"-"`
		Field4 string
	}

	t.Run("extracts field name from yaml tag", func(t *testing.T) {
		s := TestStruct{}
		v := reflect.TypeOf(s)
		field := v.Field(0)
		assert.Equal(t, "field_one", getYAMLFieldNameFromTag(field))
	})

	t.Run("extracts field name ignoring options", func(t *testing.T) {
		s := TestStruct{}
		v := reflect.TypeOf(s)
		field := v.Field(1)
		assert.Equal(t, "field_two", getYAMLFieldNameFromTag(field))
	})

	t.Run("returns empty for dash tag", func(t *testing.T) {
		s := TestStruct{}
		v := reflect.TypeOf(s)
		field := v.Field(2)
		assert.Equal(t, "", getYAMLFieldNameFromTag(field))
	})

	t.Run("returns empty for missing tag", func(t *testing.T) {
		s := TestStruct{}
		v := reflect.TypeOf(s)
		field := v.Field(3)
		assert.Equal(t, "", getYAMLFieldNameFromTag(field))
	})
}

func TestConvertToMapSlice(t *testing.T) {
	t.Run("converts slice of maps", func(t *testing.T) {
		input := []any{
			map[string]any{"key": "value1"},
			map[string]any{"key": "value2"},
		}
		result := convertToMapSlice(input)
		assert.Len(t, result, 2)
		assert.Equal(t, "value1", result[0]["key"])
		assert.Equal(t, "value2", result[1]["key"])
	})

	t.Run("skips non-map items", func(t *testing.T) {
		input := []any{
			map[string]any{"key": "value"},
			"not a map",
			123,
		}
		result := convertToMapSlice(input)
		assert.Len(t, result, 1)
		assert.Equal(t, "value", result[0]["key"])
	})

	t.Run("handles empty slice", func(t *testing.T) {
		input := []any{}
		result := convertToMapSlice(input)
		assert.Len(t, result, 0)
	})
}
