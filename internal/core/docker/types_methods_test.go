package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceCharacteristicsResolver_ResolveComposeUpFlags(t *testing.T) {
	resolver, err := NewServiceCharacteristicsResolver()
	assert.NoError(t, err)
	assert.NotNil(t, resolver)

	flags := resolver.ResolveComposeUpFlags([]string{})
	assert.NotNil(t, flags)
}

func TestServiceCharacteristicsResolver_ResolveComposeDownFlags(t *testing.T) {
	resolver, err := NewServiceCharacteristicsResolver()
	assert.NoError(t, err)
	assert.NotNil(t, resolver)

	flags := resolver.ResolveComposeDownFlags([]string{})
	assert.NotNil(t, flags)
}

func TestUpOptions_ToSDK(t *testing.T) {
	opts := UpOptions{}
	sdk := opts.ToSDK()
	assert.NotNil(t, sdk)
}

func TestDownOptions_ToSDK(t *testing.T) {
	opts := DownOptions{}
	sdk := opts.ToSDK()
	assert.NotNil(t, sdk)
}

func TestStopOptions_ToSDK(t *testing.T) {
	opts := StopOptions{}
	sdk := opts.ToSDK()
	assert.NotNil(t, sdk)
}

func TestLogOptions_ToSDK(t *testing.T) {
	opts := LogOptions{}
	sdk := opts.ToSDK()
	assert.NotNil(t, sdk)
}
