//go:build unit

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionDisplayManager_DisplayJSON(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.displayJSON("1.0.0")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayYAML(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.displayYAML("1.0.0")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayBasic_Text(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.DisplayBasic("1.0.0", "text")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayBasic_JSON(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.DisplayBasic("1.0.0", "json")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayBasic_YAML(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.DisplayBasic("1.0.0", "yaml")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayFull_Text(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.DisplayFull("1.0.0", "text")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayFull_JSON(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.DisplayFull("1.0.0", "json")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_DisplayFull_YAML(t *testing.T) {
	vdm := NewVersionDisplayManager()
	err := vdm.DisplayFull("1.0.0", "yaml")
	assert.NoError(t, err)
}

func TestVersionDisplayManager_GetCurrentVersion(t *testing.T) {
	vdm := NewVersionDisplayManager()
	version := vdm.GetCurrentVersion()
	assert.NotEmpty(t, version)
}
