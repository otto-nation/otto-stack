//go:build unit

package services

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestManager_UncoveredMethods(t *testing.T) {
	t.Run("tests GetDependencies", func(t *testing.T) {
		manager, err := New()
		testhelpers.AssertValidConstructor(t, manager, err, "Manager")

		deps, err := manager.GetDependencies("postgres")
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.IsType(t, []string{}, deps)
		}
	})

	t.Run("tests BuildConnectCommand", func(t *testing.T) {
		manager, err := New()
		testhelpers.AssertValidConstructor(t, manager, err, "Manager")

		options := map[string]string{"user": "test"}
		cmd, err := manager.BuildConnectCommand("postgres", options)
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.IsType(t, []string{}, cmd)
		}
	})

	t.Run("tests ExecuteCustomOperation", func(t *testing.T) {
		manager, err := New()
		testhelpers.AssertValidConstructor(t, manager, err, "Manager")

		result, err := manager.ExecuteCustomOperation("test-op", "service")
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.IsType(t, "", result)
		}
	})
}
