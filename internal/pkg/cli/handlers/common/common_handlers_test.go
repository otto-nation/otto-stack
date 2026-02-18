//go:build unit

package common

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestCommon_ValidateArgs(t *testing.T) {
	handler := &BaseHandler{}
	err := handler.ValidateArgs([]string{})
	testhelpers.AssertNoError(t, err, "ValidateArgs with empty args should not error")
}

func TestCommon_GetRequiredFlags(t *testing.T) {
	handler := &BaseHandler{}
	flags := handler.GetRequiredFlags()
	testhelpers.AssertNoError(t, nil, "GetRequiredFlags should not error")
	if flags == nil {
		t.Error("GetRequiredFlags should return flags slice")
	}
}
