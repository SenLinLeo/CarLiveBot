package ark

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const defaultTimeout = 60 * time.Second

// Client 方舟 Chat API 客户端
type Client struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

// NewClient 创建 Ark Chat 客户端
func NewClient(apiKey, baseURL, modelID string) *Client {
	if baseURL == "" {
		baseURL = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
	}
	if modelID == "" {
		modelID = "doubao-seed-2-0-lite-260215"
	}
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		model:   modelID,
		client:  &http.Client{Timeout: defaultTimeout},
	}
}

// Model 返回当前模型 ID
func (c *Client) Model() string { return c.model }

// ChatStream 流式请求，通过 callback 逐片返回 content；若 callback 返回 error 则中止
func (c *Client) ChatStream(ctx context.Context, req *ChatRequest, fn func(content string) error) error {
	if req.Model == "" {
		req.Model = c.model
	}
	req.Stream = true

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bs, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ark api status %d: %s", resp.StatusCode, string(bs))
	}

	return parseSSE(resp.Body, fn)
}

// parseSSE 解析 SSE 流，每行 data: {...} 解析为 StreamChunk 并回调 content
func parseSSE(r io.Reader, fn func(content string) error) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(nil, 1024*1024)
	var dataLine string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			dataLine = strings.TrimPrefix(line, "data: ")
		}
		if dataLine == "" || dataLine == "[DONE]" {
			continue
		}
		var chunk StreamChunk
		if err := json.Unmarshal([]byte(dataLine), &chunk); err != nil {
			continue
		}
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			if err := fn(chunk.Choices[0].Delta.Content); err != nil {
				return err
			}
		}
		dataLine = ""
	}
	return scanner.Err()
}
