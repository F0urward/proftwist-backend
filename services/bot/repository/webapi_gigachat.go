package repository

import (
	"context"
	"fmt"

	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	gigachatClientDTO "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient/dto"
	"github.com/F0urward/proftwist-backend/internal/server/middleware/logctx"
	"github.com/F0urward/proftwist-backend/services/bot"
)

type GigachatWebapi struct {
	client *gigachatclient.Client
}

func NewGigachatWebapi(client *gigachatclient.Client) bot.GigachatWebapi {
	return &GigachatWebapi{client: client}
}

func (r *GigachatWebapi) GetBotResponse(ctx context.Context, query string) (string, error) {
	const op = "GigachatWebapi.GetBotResponse"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	logger.WithField("query_length", len(query)).Debug("getting bot response from Gigachat")

	chatReq := &gigachatClientDTO.ChatRequest{
		Model: "GigaChat",
		Messages: []gigachatClientDTO.Message{
			{
				Role: "system",
				Content: `Ты - полезный AI-ассистент в чате образовательной платформы. 
Отвечай кратко, информативно и по делу. 
Будь вежливым и помогай пользователям с их вопросами.

ВАЖНО: Отвечай текстом без форматирования. Никих переносов строки типа \n, Markdown, эмодзи, или специальных символов. Текст в ОДНУ СТРОКУ`,
			},
			{
				Role:    "user",
				Content: query,
			},
		},
		Temperature:       float64Ptr(0.7),
		MaxTokens:         int64Ptr(1024),
		RepetitionPenalty: float64Ptr(1.1),
	}

	chatResp, err := r.client.Chat(ctx, chatReq)
	if err != nil {
		logger.WithError(err).Error("failed to get response from Gigachat")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if len(chatResp.Choices) == 0 {
		logger.Error("empty response from Gigachat")
		return "", fmt.Errorf("%s: empty response from Gigachat", op)
	}

	responseText := chatResp.Choices[0].Message.Content

	logger.WithField("response_length", len(responseText)).Debug("successfully received bot response")
	return responseText, nil
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
