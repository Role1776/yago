package yago

// ReasoningOptions holds parameters for the reasoning mode.
type ReasoningOptions struct {
	Mode string `json:"mode"`
}

// CompletionOptions holds the configuration parameters for a completion request.
type CompletionOptions struct {
	Stream           bool             `json:"stream"`
	Temperature      float64          `json:"temperature"`
	MaxTokens        string           `json:"maxTokens"`
	ReasoningOptions ReasoningOptions `json:"reasoningOptions,omitempty"`
}

// GenerativeModel represents a specific model that can be used for generation.
// It holds the configuration for API requests.
type GenerativeModel struct {
	c                 *Client
	modelURI          string
	CompletionOptions CompletionOptions
	SystemInstruction string
}

// GenerativeModel returns a new GenerativeModel instance for the specified folderID and model name.
func (c *Client) GenerativeModel(name string) *GenerativeModel {
	return &GenerativeModel{
		c:        c,
		modelURI: "gpt://" + c.folderID + "/" + name,
		CompletionOptions: CompletionOptions{
			Stream:      false,
			Temperature: 0.5,
			MaxTokens:   "9999999",
			ReasoningOptions: ReasoningOptions{
				Mode: "DISABLED",
			},
		},
		SystemInstruction: "",
	}
}
