package stack

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildStackContext(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]interface{}
		expected struct {
			projectName string
			serviceLen  int
			forceFlag   bool
		}
	}{
		{
			name: "empty args and flags",
			args: []string{},
			flags: map[string]interface{}{
				"force": false,
			},
			expected: struct {
				projectName string
				serviceLen  int
				forceFlag   bool
			}{
				projectName: "default-project",
				serviceLen:  0,
				forceFlag:   false,
			},
		},
		{
			name: "with force flag",
			args: []string{},
			flags: map[string]interface{}{
				"force": true,
			},
			expected: struct {
				projectName string
				serviceLen  int
				forceFlag   bool
			}{
				projectName: "default-project",
				serviceLen:  0,
				forceFlag:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock command with flags
			cmd := &cobra.Command{}
			cmd.Flags().Bool("force", false, "force flag")

			// Set flag values
			for key, value := range tt.flags {
				if boolVal, ok := value.(bool); ok {
					if boolVal {
						err := cmd.Flags().Set(key, "true")
						require.NoError(t, err)
					} else {
						err := cmd.Flags().Set(key, "false")
						require.NoError(t, err)
					}
				}
			}

			// Build context
			ctx, err := BuildStackContext(cmd, tt.args)
			require.NoError(t, err)

			// Verify context
			assert.Equal(t, tt.expected.projectName, ctx.Project.Name)
			assert.Equal(t, tt.expected.serviceLen, len(ctx.Services.Names))
			assert.Equal(t, tt.expected.forceFlag, ctx.Runtime.Force)
		})
	}
}
