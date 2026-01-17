package output

// ToolOutput represents the structured output from a statusline tool
type ToolOutput struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
