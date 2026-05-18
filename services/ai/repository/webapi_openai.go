package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/ai"
	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

const defaultOpenAIBaseURL = "https://api.openai.com/v1"
const defaultOllamaBaseURL = "http://localhost:11434/v1"

type OpenAICompatibleProvider struct {
	baseURL       string
	apiKey        string
	model         string
	requireAPIKey bool
	client        *http.Client
}

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message openAIResponseMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

type openAIResponseMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type openAIContentPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func NewOpenAICompatibleProviderWithCredentials(baseURL, apiKey, model string) ai.Provider {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = defaultOpenAIBaseURL
	}

	return &OpenAICompatibleProvider{
		baseURL:       baseURL,
		apiKey:        strings.TrimSpace(apiKey),
		model:         strings.TrimSpace(model),
		requireAPIKey: true,
		client:        http.DefaultClient,
	}
}

func NewOllamaProviderWithCredentials(baseURL, apiKey, model string) ai.Provider {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = defaultOllamaBaseURL
	}

	return &OpenAICompatibleProvider{
		baseURL:       baseURL,
		apiKey:        strings.TrimSpace(apiKey),
		model:         strings.TrimSpace(model),
		requireAPIKey: false,
		client:        http.DefaultClient,
	}
}

func (p *OpenAICompatibleProvider) GenerateRoadmapNodeDescription(ctx context.Context, req dto.GenerateRoadmapNodeDescriptionRequestDTO) (string, error) {
	const op = "OpenAICompatibleProvider.GenerateRoadmapNodeDescription"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	if p.requireAPIKey && p.apiKey == "" {
		return "", fmt.Errorf("%s: AI OpenAI-compatible API key is not configured", op)
	}
	if p.model == "" {
		return "", fmt.Errorf("%s: AI OpenAI-compatible model is not configured", op)
	}

	chatReq := openAIChatRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{
				Role:    "user",
				Content: req.NodeLabel,
			},
		},
		Temperature: 0.4,
		MaxTokens:   350,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("%s: failed to marshal chat request: %w", op, err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("%s: failed to create chat request: %w", op, err)
	}
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("%s: request canceled while waiting for provider: %w", op, ctx.Err())
		}
		logger.WithError(err).Error("failed to send OpenAI-compatible request")
		return "", fmt.Errorf("%s: failed to send chat request: %w", op, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("%s: request canceled while reading provider response: %w", op, ctx.Err())
		}
		return "", fmt.Errorf("%s: failed to read chat response: %w", op, err)
	}

	var chatResp openAIChatResponse
	if err = json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("%s: failed to decode chat response: %w", op, err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		message := resp.Status
		if chatResp.Error != nil && chatResp.Error.Message != "" {
			message = chatResp.Error.Message
		}
		return "", fmt.Errorf("%s: provider returned error: %s", op, message)
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("%s: empty response from provider", op)
	}

	description := strings.TrimSpace(extractOpenAIMessageContent(chatResp.Choices[0].Message.Content))
	if description == "" {
		return "", fmt.Errorf("%s: empty description from provider", op)
	}

	return description, nil
}

func (p *OpenAICompatibleProvider) GenerateRoadmap(ctx context.Context, req dto.GenerateRoadmapRequestDTO) (string, error) {
	const op = "OpenAICompatibleProvider.GenerateRoadmap"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	if p.requireAPIKey && p.apiKey == "" {
		return "", fmt.Errorf("%s: AI OpenAI-compatible API key is not configured", op)
	}
	if p.model == "" {
		return "", fmt.Errorf("%s: AI OpenAI-compatible model is not configured", op)
	}

	chatReq := openAIChatRequest{
		Model: p.model,
		Messages: []openAIMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		Temperature: 0.35,
		MaxTokens:   6000,
	}

	respBody, err := p.sendChatRequest(ctx, chatReq)
	if err != nil {
		logger.WithError(err).Error("failed to send OpenAI-compatible roadmap request")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return respBody, nil
}

func (p *OpenAICompatibleProvider) sendChatRequest(ctx context.Context, chatReq openAIChatRequest) (string, error) {
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create chat request: %w", err)
	}
	if p.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("request canceled while waiting for provider: %w", ctx.Err())
		}
		return "", fmt.Errorf("failed to send chat request: %w", err)
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("request canceled while reading provider response: %w", ctx.Err())
		}
		return "", fmt.Errorf("failed to read chat response: %w", err)
	}

	var chatResp openAIChatResponse
	if err = json.Unmarshal(rawBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to decode chat response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		message := resp.Status
		if chatResp.Error != nil && chatResp.Error.Message != "" {
			message = chatResp.Error.Message
		}
		return "", fmt.Errorf("provider returned error: %s", message)
	}
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from provider")
	}

	content := strings.TrimSpace(extractOpenAIMessageContent(chatResp.Choices[0].Message.Content))
	if content == "" {
		return "", fmt.Errorf("empty content from provider")
	}

	return content, nil
}

func extractOpenAIMessageContent(raw json.RawMessage) string {
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		return text
	}

	var parts []openAIContentPart
	if err := json.Unmarshal(raw, &parts); err == nil {
		var b strings.Builder
		for _, part := range parts {
			if part.Text == "" {
				continue
			}
			if b.Len() > 0 {
				b.WriteString("\n")
			}
			b.WriteString(part.Text)
		}
		return b.String()
	}

	return ""
}
