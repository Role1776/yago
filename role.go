package yago

// Role represents the role of a message's author in a chat conversation.
// It's used to distinguish between system instructions, user inputs, and assistant responses.
type Role string

const (
	// RoleUser specifies the role of the end-user providing input.
	RoleUser Role = "user"

	// RoleSystem specifies a system-level instruction or context for the model.
	RoleSystem Role = "system"

	// RoleAssistant specifies the role of the AI model generating the response.
	RoleAssistant Role = "assistant"
)

// Message represents a single message in a chat sequence.
// It contains the role of the author and the text content of the message.
type Message struct {
	Role Role   `json:"role"`
	Text string `json:"text"`
}
