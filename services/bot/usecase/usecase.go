package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/config"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/bot"
	"github.com/F0urward/proftwist-backend/services/bot/dto"
)

type BotConfig struct {
	botUserID        string
	botTriggerPhrase string
}

type BotUsecase struct {
	gigachatWebapi bot.GigachatWebapi
	chatClient     chatclient.ChatServiceClient
	botConfig      BotConfig
}

func NewBotUsecase(
	gigachatWebapi bot.GigachatWebapi,
	chatClient chatclient.ChatServiceClient,
	cfg *config.Config,
) bot.Usecase {
	return &BotUsecase{
		gigachatWebapi: gigachatWebapi,
		chatClient:     chatClient,
		botConfig: BotConfig{
			botUserID:        cfg.Bot.BotUserID,
			botTriggerPhrase: cfg.Bot.BotTriggerPhrase,
		},
	}
}

func (uc *BotUsecase) HandleBotTrigger(ctx context.Context, event dto.MessageForBotEvent) error {
	const op = "BotUsecase.HandleBotTrigger"
	logger := ctxutil.GetLogger(ctx).WithField("op", op).WithField("chat_id", event.ChatID)

	logger.Info("handling bot trigger message")

	if !uc.isBotTrigger(event.Content) {
		logger.Info("message is not a bot trigger, ignoring")
		return nil
	}

	query := uc.extractQuery(event.Content)
	if query == "" {
		logger.Warn("empty query after removing trigger phrase")
		return nil
	}

	logger.WithField("query", query).Info("processing query with Gigachat")

	response, err := uc.gigachatWebapi.GetBotResponse(ctx, query, event.ChatTitle)
	if err != nil {
		logger.WithError(err).Error("failed to get bot response from Gigachat")
		return fmt.Errorf("%s: %w", op, err)
	}

	cleanedResponse := uc.cleanBotResponse(response)

	sendReq := &chatclient.SendGroupChatMessageRequest{
		ChatId:  event.ChatID,
		UserId:  uc.botConfig.botUserID,
		Content: cleanedResponse,
	}

	_, err = uc.chatClient.SendGroupChatMessage(ctx, sendReq)
	if err != nil {
		logger.WithError(err).Error("failed to send bot response via chat client")
		return fmt.Errorf("%s: %w", op, err)
	}

	logger.WithField("response_length", len(cleanedResponse)).Info("successfully handled bot trigger")
	return nil
}

func (uc *BotUsecase) isBotTrigger(content string) bool {
	return strings.HasPrefix(strings.TrimSpace(content), uc.botConfig.botTriggerPhrase)
}

func (uc *BotUsecase) extractQuery(content string) string {
	if strings.HasPrefix(strings.TrimSpace(content), uc.botConfig.botTriggerPhrase) {
		return strings.TrimSpace(content[len(uc.botConfig.botTriggerPhrase):])
	}
	return strings.TrimSpace(content)
}

func (uc *BotUsecase) cleanBotResponse(response string) string {
	trimmedResponse := strings.TrimSpace(response)

	if strings.HasPrefix(strings.ToLower(trimmedResponse), strings.ToLower(uc.botConfig.botTriggerPhrase)) {
		cleaned := strings.TrimSpace(trimmedResponse[len(uc.botConfig.botTriggerPhrase):])
		if cleaned == "" {
			return "К сожалению, я не знаю, что ответить, на ваш вопрос."
		}
		return cleaned
	}

	return response
}
