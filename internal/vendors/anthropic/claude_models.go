package anthropic

import (
	"github.com/baalimago/clai/internal/tools"
)

type ClaudeResponse struct {
	Content      []ClaudeMessage `json:"content"`
	ID           string          `json:"id"`
	Model        string          `json:"model"`
	Role         string          `json:"role"`
	StopReason   string          `json:"stop_reason"`
	StopSequence any             `json:"stop_sequence"`
	Type         string          `json:"type"`
	Usage        TokenInfo       `json:"usage"`
}

type ClaudeMessage struct {
	ID    string      `json:"id,omitempty"`
	Input tools.Input `json:"input,omitempty"`
	Name  string      `json:"name,omitempty"`
	Text  string      `json:"text,omitempty"`
	Type  string      `json:"type"`
}

type TokenInfo struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type Delta struct {
	Type        string `json:"type"`
	Text        string `json:"text,omitempty"`
	PartialJson string `json:"partial_json,omitempty"`
}

type ContentBlockDelta struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta Delta  `json:"delta"`
}

type ContentBlockSuper struct {
	Type         string       `json:"type"`
	Index        int          `json:"index"`
	ContentBlock ContentBlock `json:"content_block"`
}

type ContentBlock struct {
	Type  string                 `json:"type"`
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

type Root struct {
	Type         string       `json:"type"`
	Index        int          `json:"index"`
	ContentBlock ContentBlock `json:"content_block"`
}
