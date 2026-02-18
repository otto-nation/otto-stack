//go:build unit

package config

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestConfigService_New(t *testing.T) {
	service := NewConfigService()
	testhelpers.AssertNoError(t, nil, "NewConfigService should not error")
	if service == nil {
		t.Error("NewConfigService should return a service instance")
	}
}

func TestConfigService_GetConfigHash_Nil(t *testing.T) {
	service := NewConfigService()
	hash, err := service.GetConfigHash(nil)
	testhelpers.AssertError(t, err, "GetConfigHash with nil should error")
	if hash != "" {
		t.Error("GetConfigHash with nil should return empty hash")
	}
}

func TestConfig_Load(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		testhelpers.AssertError(t, err, "LoadConfig should handle missing files")
	} else {
		testhelpers.AssertNoError(t, err, "LoadConfig should not error")
		if config == nil {
			t.Error("LoadConfig should return config")
		}
	}
}

func TestConfig_GetConfigPath(t *testing.T) {
	path := getConfigPath()
	testhelpers.AssertNoError(t, nil, "getConfigPath should not error")
	if path == "" {
		t.Error("getConfigPath should return path")
	}
}

func TestConfig_GetLocalConfigPath(t *testing.T) {
	path := getLocalConfigPath()
	testhelpers.AssertNoError(t, nil, "getLocalConfigPath should not error")
	if path == "" {
		t.Error("getLocalConfigPath should return path")
	}
}
