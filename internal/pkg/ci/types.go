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

// StatusOutput represents service status output
type StatusOutput struct {
	Services []any `json:"services"`
	Count    int   `json:"count"`
}

// InterfacesOutput represents web interfaces output
type InterfacesOutput struct {
	Interfaces []any `json:"interfaces"`
	Count      int   `json:"count"`
}
