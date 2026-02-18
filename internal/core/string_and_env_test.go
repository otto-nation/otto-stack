//go:build unit

package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTitleCase_Hyphenated(t *testing.T) {
	result := TitleCase("hello-world")
	assert.Equal(t, "Hello World", result)
}

func TestTitleCase_SingleWord(t *testing.T) {
	result := TitleCase("hello")
	assert.Equal(t, "Hello", result)
}

func TestTitleCase_Empty(t *testing.T) {
	result := TitleCase("")
	assert.Equal(t, "", result)
}

func TestTitleCase_MultipleHyphens(t *testing.T) {
	result := TitleCase("hello-world-test")
	assert.Equal(t, "Hello World Test", result)
}

func TestTitleCase_NoHyphens(t *testing.T) {
	result := TitleCase("hello")
	assert.Equal(t, "Hello", result)
}
