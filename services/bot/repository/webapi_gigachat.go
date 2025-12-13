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

func (r *GigachatWebapi) GetBotResponse(ctx context.Context, query, chatTitle string) (string, error) {
	const op = "GigachatWebapi.GetBotResponse"
	logger := logctx.GetLogger(ctx).WithField("op", op)

	logger.WithFields(map[string]interface{}{
		"query_length": len(query),
		"chat_title":   chatTitle,
	}).Info("getting bot response from Gigachat")

	systemPrompt := `Ты - полезный AI-ассистент в чате образовательной платформы. 
Тебе создала команда Fourward.
Отвечай кратко, информативно и по делу. 
Будь вежливым и помогай пользователям с их вопросами.

Текущий чат: "` + chatTitle + `"
ВАЖНО: Отвечай текстом без форматирования. Никаких переносов строки типа \n, Markdown, эмодзи, или специальных символов. Текст в ОДНУ СТРОКУ`

	chatReq := &gigachatClientDTO.ChatRequest{
		Model: "GigaChat",
		Messages: []gigachatClientDTO.Message{
			{
				Role:    "system",
				Content: systemPrompt,
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

	logger.WithField("response_length", len(responseText)).Info("successfully received bot response")
	return responseText, nil
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
