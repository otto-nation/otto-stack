//go:build unit

package version

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/otto-nation/otto-stack/test/testhelpers"
)

func TestCheckForUpdates_NewerRelease(t *testing.T) {
	if IsDevBuild() {
		t.Skip("Skipping test in dev build")
	}

	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				body := `{
					"tag_name": "v2.0.0",
					"name": "Release 2.0.0",
					"body": "New features",
					"draft": false,
					"prerelease": false,
					"html_url": "https://github.com/test/repo/releases/tag/v2.0.0"
				}`
				return testhelpers.MockJSONResponse(200, body), nil
			},
		},
	}

	release, hasUpdate, err := checker.CheckForUpdates()
	require.NoError(t, err)
	assert.True(t, hasUpdate)
	assert.NotNil(t, release)
	assert.Equal(t, "v2.0.0", release.TagName)
}

func TestCheckForUpdates_SameVersion(t *testing.T) {
	if IsDevBuild() {
		t.Skip("Skipping test in dev build")
	}

	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				body := `{
					"tag_name": "v1.0.0",
					"name": "Release 1.0.0",
					"draft": false,
					"prerelease": false
				}`
				return testhelpers.MockJSONResponse(200, body), nil
			},
		},
	}

	release, hasUpdate, err := checker.CheckForUpdates()
	require.NoError(t, err)
	assert.False(t, hasUpdate)
	assert.NotNil(t, release)
}

func TestCheckForUpdates_SkipsDraft(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				body := `{
					"tag_name": "v2.0.0",
					"draft": true,
					"prerelease": false
				}`
				return testhelpers.MockJSONResponse(200, body), nil
			},
		},
	}

	release, hasUpdate, err := checker.CheckForUpdates()
	require.NoError(t, err)
	assert.False(t, hasUpdate)
	assert.Nil(t, release)
}

func TestCheckForUpdates_SkipsPrerelease(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				body := `{
					"tag_name": "v2.0.0-beta",
					"draft": false,
					"prerelease": true
				}`
				return testhelpers.MockJSONResponse(200, body), nil
			},
		},
	}

	release, hasUpdate, err := checker.CheckForUpdates()
	require.NoError(t, err)
	assert.False(t, hasUpdate)
	assert.Nil(t, release)
}

func TestCheckForUpdates_HTTPError(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		},
	}

	_, _, err := checker.CheckForUpdates()
	assert.Error(t, err)
}

func TestCheckForUpdates_Non200Status(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return testhelpers.MockResponse(404, "Not Found"), nil
			},
		},
	}

	_, _, err := checker.CheckForUpdates()
	assert.Error(t, err)
}

func TestCheckForUpdates_InvalidJSON(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return testhelpers.MockJSONResponse(200, "invalid json"), nil
			},
		},
	}

	_, _, err := checker.CheckForUpdates()
	assert.Error(t, err)
}

func TestCheckForUpdates_Headers(t *testing.T) {
	checker := NewUpdateChecker("1.0.0")
	var capturedReq *http.Request
	checker.client = &http.Client{
		Transport: &mockTransport{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				body := `{"tag_name": "v1.0.0", "draft": false, "prerelease": false}`
				return testhelpers.MockJSONResponse(200, body), nil
			},
		},
	}

	_, _, _ = checker.CheckForUpdates()
	assert.NotNil(t, capturedReq)
	assert.NotEmpty(t, capturedReq.Header.Get("User-Agent"))
	assert.Equal(t, "application/vnd.github.v3+json", capturedReq.Header.Get("Accept"))
}

// mockTransport implements http.RoundTripper for testing
type mockTransport struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
