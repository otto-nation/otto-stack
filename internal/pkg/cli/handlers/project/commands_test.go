package project

import (
	"testing"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestProjectCommands_constructors(t *testing.T) {
	t.Run("new conflicts command", func(t *testing.T) {
		cmd := NewConflictsCommand()
		testhelpers.AssertNoError(t, nil, "NewConflictsCommand should not error")
		if cmd == nil {
			t.Error("NewConflictsCommand should return a command")
		}
	})

	t.Run("new validate command", func(t *testing.T) {
		cmd := NewValidateCommand()
		testhelpers.AssertNoError(t, nil, "NewValidateCommand should not error")
		if cmd == nil {
			t.Error("NewValidateCommand should return a command")
		}
	})

	t.Run("new doctor command", func(t *testing.T) {
		cmd := NewDoctorCommand()
		testhelpers.AssertNoError(t, nil, "NewDoctorCommand should not error")
		if cmd == nil {
			t.Error("NewDoctorCommand should return a command")
		}
	})
}
