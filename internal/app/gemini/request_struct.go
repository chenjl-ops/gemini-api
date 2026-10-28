package gemini

type AiRequest struct {
	Prompt string `json:"prompt"`
	ApiKey string `json:"api_key"`
}
