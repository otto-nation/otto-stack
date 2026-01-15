//go:build unit

package common

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestCommon_base_functions(t *testing.T) {
	handler := &BaseHandler{}

	t.Run("validate args with empty args", func(t *testing.T) {
		err := handler.ValidateArgs([]string{})
		testhelpers.AssertNoError(t, err, "ValidateArgs with empty args should not error")
	})

	t.Run("get required flags", func(t *testing.T) {
		flags := handler.GetRequiredFlags()
		testhelpers.AssertNoError(t, nil, "GetRequiredFlags should not error")
		if flags == nil {
			t.Error("GetRequiredFlags should return flags slice")
		}
	})
}
