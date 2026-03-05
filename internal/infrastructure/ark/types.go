package ark

// ChatMessage 单条消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 请求体（流式）
type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      bool          `json:"stream"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Thinking    *Thinking     `json:"thinking,omitempty"`
}

// Thinking 深度思考
type Thinking struct {
	Type string `json:"type"` // "enabled"
}

// SSE 流式 chunk（只解析 content）
type StreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}
