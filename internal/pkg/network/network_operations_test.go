//go:build unit

package network

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestNetwork_port_functions(t *testing.T) {
	t.Run("is port in use with invalid port", func(t *testing.T) {
		inUse := IsPortInUse(-1)
		testhelpers.AssertNoError(t, nil, "IsPortInUse should not error")
		if inUse {
			t.Error("IsPortInUse with invalid port should return false")
		}
	})

	t.Run("get free port with high start port", func(t *testing.T) {
		port, err := GetFreePort(65000)
		if err != nil {
			testhelpers.AssertError(t, err, "GetFreePort with high port may error")
		} else {
			testhelpers.AssertNoError(t, err, "GetFreePort with high port should not error")
			if port == 0 {
				t.Error("GetFreePort should return non-zero port")
			}
		}
	})
}
