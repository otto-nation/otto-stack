//go:build unit

package services

import (
	"reflect"
	"testing"

	"github.com/otto-nation/otto-stack/internal/core/docker"
	"github.com/otto-nation/otto-stack/internal/pkg/logger"
	"github.com/otto-nation/otto-stack/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestHasInitScripts(t *testing.T) {
	s := &Service{logger: logger.GetLogger()}

	config := fixtures.NewServiceConfig("test").Build()
	config.InitService = &docker.InitServiceSpec{Enabled: true}
	assert.True(t, s.hasInitScripts(config))

	config.InitService.Enabled = false
	assert.False(t, s.hasInitScripts(config))

	config.InitService = nil
	assert.False(t, s.hasInitScripts(config))
}

func TestHasLocalInitScripts(t *testing.T) {
	s := &Service{logger: logger.GetLogger()}

	config := fixtures.NewServiceConfig("test").Build()
	config.InitService = &docker.InitServiceSpec{Enabled: true, Mode: docker.InitServiceModeLocal}
	assert.True(t, s.hasLocalInitScripts(config))

	config.InitService.Mode = docker.InitServiceModeContainer
	assert.False(t, s.hasLocalInitScripts(config))

	config.InitService.Enabled = false
	config.InitService.Mode = docker.InitServiceModeLocal
	assert.False(t, s.hasLocalInitScripts(config))
}

func TestIsConfigStructField(t *testing.T) {
	type TestConfig struct {
		Value string
	}

	type TestStruct struct {
		Config *TestConfig
		Name   string
	}

	s := TestStruct{}
	v := reflect.ValueOf(&s).Elem()
	field := v.Field(0)
	fieldType := v.Type().Field(0)
	assert.True(t, isConfigStructField(field, fieldType))

	field = v.Field(1)
	fieldType = v.Type().Field(1)
	assert.False(t, isConfigStructField(field, fieldType))
}

func TestGetYAMLFieldNameFromTag(t *testing.T) {
	type TestStruct struct {
		Field1 string `yaml:"field_one"`
		Field2 string `yaml:"field_two,omitempty"`
		Field3 string `yaml:"-"`
		Field4 string
	}

	s := TestStruct{}
	v := reflect.TypeOf(s)
	field := v.Field(0)
	assert.Equal(t, "field_one", getYAMLFieldNameFromTag(field))

	field = v.Field(1)
	assert.Equal(t, "field_two", getYAMLFieldNameFromTag(field))

	field = v.Field(2)
	assert.Equal(t, "", getYAMLFieldNameFromTag(field))

	field = v.Field(3)
	assert.Equal(t, "", getYAMLFieldNameFromTag(field))
}

func TestConvertToMapSlice(t *testing.T) {
	input := []any{
		map[string]any{"key": "value1"},
		map[string]any{"key": "value2"},
	}
	result := convertToMapSlice(input)
	assert.Len(t, result, 2)
	assert.Equal(t, "value1", result[0]["key"])
	assert.Equal(t, "value2", result[1]["key"])

	input = []any{
		map[string]any{"key": "value"},
		"not a map",
		123,
	}
	result = convertToMapSlice(input)
	assert.Len(t, result, 1)
	assert.Equal(t, "value", result[0]["key"])

	input = []any{}
	result = convertToMapSlice(input)
	assert.Len(t, result, 0)
}
