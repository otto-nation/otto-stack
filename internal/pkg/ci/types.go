package ci

// Output represents structured command output
type Output struct {
	Result   any `json:"result"`
	ExitCode int `json:"exit_code"`
}

// ErrorOutput represents error output
type ErrorOutput struct {
	Error    string `json:"error"`
	ExitCode int    `json:"exit_code"`
}
