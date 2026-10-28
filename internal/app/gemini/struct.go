package gemini

// InteractionRequest Interactions API 请求体
type InteractionRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// StepContent 2. Interactions API 响应体结构
type StepContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InteractionStep struct {
	Type    string        `json:"type"`
	Content []StepContent `json:"content,omitempty"`
}

type InteractionResponse struct {
	ID         string            `json:"id"`
	Status     string            `json:"status"`
	Model      string            `json:"model"`
	OutputText string            `json:"output_text,omitempty"`
	Steps      []InteractionStep `json:"steps,omitempty"`
	Error      *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}
