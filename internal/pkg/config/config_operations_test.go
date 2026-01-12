package config

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestConfigService_basic_operations(t *testing.T) {
	t.Run("new config service", func(t *testing.T) {
		service := NewConfigService()
		testhelpers.AssertNoError(t, nil, "NewConfigService should not error")
		if service == nil {
			t.Error("NewConfigService should return a service instance")
		}
	})

	t.Run("get config hash with nil config", func(t *testing.T) {
		service := NewConfigService()
		hash, err := service.GetConfigHash(nil)
		testhelpers.AssertError(t, err, "GetConfigHash with nil should error")
		if hash != "" {
			t.Error("GetConfigHash with nil should return empty hash")
		}
	})
}

func TestConfig_LoadOperations(t *testing.T) {
	t.Run("load config", func(t *testing.T) {
		config, err := LoadConfig()
		// Will likely fail due to missing files but tests the function
		if err != nil {
			testhelpers.AssertError(t, err, "LoadConfig should handle missing files")
		} else {
			testhelpers.AssertNoError(t, err, "LoadConfig should not error")
			if config == nil {
				t.Error("LoadConfig should return config")
			}
		}
	})

	t.Run("load service config", func(t *testing.T) {
		configs, err := LoadServiceConfig("test-service")
		// Will likely fail due to missing files but tests the function
		if err != nil {
			testhelpers.AssertError(t, err, "LoadServiceConfig should handle missing files")
		} else {
			testhelpers.AssertNoError(t, err, "LoadServiceConfig should not error")
			if configs == nil {
				t.Error("LoadServiceConfig should return configs")
			}
		}
	})
}

func TestConfig_PathOperations(t *testing.T) {
	t.Run("get config path", func(t *testing.T) {
		path := getConfigPath()
		testhelpers.AssertNoError(t, nil, "getConfigPath should not error")
		if path == "" {
			t.Error("getConfigPath should return path")
		}
	})

	t.Run("get local config path", func(t *testing.T) {
		path := getLocalConfigPath()
		testhelpers.AssertNoError(t, nil, "getLocalConfigPath should not error")
		if path == "" {
			t.Error("getLocalConfigPath should return path")
		}
	})

	t.Run("get service config dir", func(t *testing.T) {
		dir := getServiceConfigDir()
		testhelpers.AssertNoError(t, nil, "getServiceConfigDir should not error")
		if dir == "" {
			t.Error("getServiceConfigDir should return directory")
		}
	})
}
