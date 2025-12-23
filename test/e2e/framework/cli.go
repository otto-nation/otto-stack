package framework

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type CLIRunner struct {
	t       *testing.T
	binPath string
	workDir string
	envVars map[string]string
}

type CLIResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

func NewCLIRunner(t *testing.T, binPath, workDir string) *CLIRunner {
	return &CLIRunner{
		t:       t,
		binPath: binPath,
		workDir: workDir,
		envVars: make(map[string]string),
	}
}

func (c *CLIRunner) SetEnv(key, value string) {
	c.envVars[key] = value
}

func (c *CLIRunner) WorkDir() string {
	return c.workDir
}

func (c *CLIRunner) Run(args ...string) *CLIResult {
	c.t.Helper()

	cmd := exec.Command(c.binPath, args...)
	cmd.Dir = c.workDir

	// Inherit parent environment and add custom vars
	cmd.Env = os.Environ()
	for k, v := range c.envVars {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	stdout, err := cmd.Output()
	result := &CLIResult{
		Stdout: string(stdout),
		Error:  err,
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		result.Stderr = string(exitError.Stderr)
		result.ExitCode = exitError.ExitCode()
	}

	return result
}

func (c *CLIRunner) RunExpectSuccess(args ...string) *CLIResult {
	c.t.Helper()

	result := c.Run(args...)
	if result.Error != nil {
		c.t.Logf("Command stderr: %s", result.Stderr)
		c.t.Logf("Command stdout: %s", result.Stdout)
	}
	require.NoError(c.t, result.Error, "Command failed: %s", strings.Join(args, " "))
	return result
}

func (c *CLIRunner) RunExpectSuccessWithBuilder(builder func() []string) *CLIResult {
	return c.RunExpectSuccess(builder()...)
}

func (c *CLIRunner) RunWithBuilder(builder func() []string) *CLIResult {
	return c.Run(builder()...)
}

type BinaryBuilder struct {
	t *testing.T
}

func NewBinaryBuilder(t *testing.T) *BinaryBuilder {
	return &BinaryBuilder{t: t}
}

func (b *BinaryBuilder) Build(outputPath string) string {
	b.t.Helper()

	// Build from the project root - go up from test/e2e to project root
	cmd := exec.Command("go", "build", "-o", outputPath, "./cmd/otto-stack")
	cmd.Dir = "../../" // Go up from test/e2e to project root
	output, err := cmd.CombinedOutput()
	if err != nil {
		b.t.Logf("Build output: %s", string(output))
	}
	require.NoError(b.t, err, "Failed to build binary")

	return outputPath
}
