package repository

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/tidwall/gjson"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	gigachatClientDTO "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient/dto"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/moderation"
)

//go:embed prompts/*
var moderationPrompts embed.FS

type ModerationGigaChatWebapi struct {
	client *gigachatclient.Client
}

func NewModerationGigaChatWebapi(client *gigachatclient.Client) moderation.GigachatWebapi {
	return &ModerationGigaChatWebapi{client: client}
}

func (r *ModerationGigaChatWebapi) GetModerationResult(ctx context.Context, content string) (*entities.ModerationResult, error) {
	const op = "GigaChatWebapi.GetModerationResult"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	logger.Info("analyzing content moderation")

	prompt, err := r.buildModerationPrompt(content)
	if err != nil {
		logger.WithError(err).Error("failed to build moderation prompt")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatReq := &gigachatClientDTO.ChatRequest{
		Model: "GigaChat",
		Messages: []gigachatClientDTO.Message{
			{
				Role:    "system",
				Content: "Ты - система модерации контента. Ты должен возвращать ТОЛЬКО валидный JSON без каких-либо дополнительных комментариев, текста или разметки.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature:       float64Ptr(0.1),
		MaxTokens:         int64Ptr(500),
		RepetitionPenalty: float64Ptr(1.0),
	}

	logger.Info("sending moderation request to GigaChat")
	chatResp, err := r.client.Chat(ctx, chatReq)
	if err != nil {
		logger.WithError(err).Error("failed to get moderation response from GigaChat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	responseText := chatResp.Choices[0].Message.Content
	result, err := r.parseModerationResponse(responseText)
	if err != nil {
		logger.WithError(err).Error("failed to parse moderation response")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.WithFields(map[string]interface{}{
		"allowed":    result.Allowed,
		"categories": result.Categories,
	}).Info("moderation result")

	return result, nil
}

func (r *ModerationGigaChatWebapi) buildModerationPrompt(content string) (string, error) {
	promptTmpl, err := moderationPrompts.ReadFile("prompts/moderation_prompt.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read moderation prompt template: %w", err)
	}

	tmpl, err := template.New("moderationPrompt").Parse(string(promptTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse moderation prompt template: %w", err)
	}

	type ModerationData struct {
		Content string
	}

	data := ModerationData{
		Content: content,
	}

	buf := &strings.Builder{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute moderation prompt template: %w", err)
	}

	return buf.String(), nil
}

func (r *ModerationGigaChatWebapi) parseModerationResponse(responseText string) (*entities.ModerationResult, error) {
	jsonData, err := r.extractJSON(responseText)
	if err != nil {
		return nil, fmt.Errorf("failed to extract JSON from moderation response: %w", err)
	}

	var result entities.ModerationResult
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal moderation result: %w", err)
	}

	return &result, nil
}

func (r *ModerationGigaChatWebapi) extractJSON(text string) (string, error) {
	result := gjson.Parse(text)

	if result.Exists() && result.Type != gjson.Null {
		return result.String(), nil
	}

	return "", fmt.Errorf("valid JSON not found in response")
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
