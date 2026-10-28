package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Text 提取 AI 回复内容的便捷方法
func (r *InteractionResponse) Text() string {
	if r.OutputText != "" {
		return r.OutputText
	}
	for _, s := range r.Steps {
		if s.Type == "model_output" {
			for _, c := range s.Content {
				if c.Text != "" {
					return c.Text
				}
			}
		}
	}
	return ""
}

// CallGeminiInteraction 核心方法
func CallGeminiInteraction(apiKey, prompt string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY 不能为空")
	}

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	apiURL := "https://generativelanguage.googleapis.com/v1beta/interactions"

	// 使用 InteractionRequest 结构体
	reqPayload := InteractionRequest{
		Model: "gemini-3-flash-preview",
		Input: prompt,
	}

	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("JSON 序列化失败: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建 Request 失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API 状态码 [%d]: %s", resp.StatusCode, string(body))
	}

	var result InteractionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("业务报错: %s", result.Error.Message)
	}

	return result.Text(), nil
}
