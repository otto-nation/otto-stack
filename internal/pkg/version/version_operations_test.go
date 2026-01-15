package version

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestVersion_validation_functions(t *testing.T) {
	t.Run("validate constraint with invalid constraint", func(t *testing.T) {
		err := ValidateConstraint("invalid")
		testhelpers.AssertError(t, err, "ValidateConstraint with invalid constraint should error")
	})

	t.Run("validate constraint with valid constraint", func(t *testing.T) {
		err := ValidateConstraint(">=1.0.0")
		testhelpers.AssertNoError(t, err, "ValidateConstraint with valid constraint should not error")
	})
}
