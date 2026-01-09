package repository

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/F0urward/proftwist-backend/internal/entities"
	"github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient"
	gigachatClientDTO "github.com/F0urward/proftwist-backend/internal/infrastructure/client/gigachatclient/dto"
	"github.com/F0urward/proftwist-backend/pkg/ctxutil"
	"github.com/F0urward/proftwist-backend/services/roadmap"
	"github.com/F0urward/proftwist-backend/services/roadmap/dto"
)

//go:embed prompts/*
var roadmapPrompts embed.FS

type RoadmapGigaChatWebapi struct {
	client *gigachatclient.Client
}

func NewRoadmapGigaChatWebapi(client *gigachatclient.Client) roadmap.GigachatWebapi {
	return &RoadmapGigaChatWebapi{client: client}
}

func (r *RoadmapGigaChatWebapi) GenerateRoadmapContent(ctx context.Context, req *dto.GenerateRoadmapDTO) (*entities.Roadmap, error) {
	const op = "GigaChatWebapi.GenerateRoadmapContent"
	logger := ctxutil.GetLogger(ctx).WithField("op", op)

	logger.WithFields(map[string]interface{}{
		"topic":      req.Topic,
		"complexity": req.Complexity,
	}).Info("generating roadmap content with GigaChat")

	example, err := r.loadRoadmapExample()
	if err != nil {
		logger.WithError(err).Error("failed to load roadmap example")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	prompt, err := r.buildGenerationPrompt(req, example)
	if err != nil {
		logger.WithError(err).Error("failed to build roadmap prompt")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	chatReq := &gigachatClientDTO.ChatRequest{
		Model: "GigaChat",
		Messages: []gigachatClientDTO.Message{
			{
				Role:    "system",
				Content: "Ты - эксперт по созданию образовательных roadmap. Ты должен возвращать ТОЛЬКО валидный JSON без каких-либо дополнительных комментариев, текста или разметки. Все ID узлов должны быть в формате UUID v4.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature:       float64Ptr(0.3),
		MaxTokens:         int64Ptr(6000),
		RepetitionPenalty: float64Ptr(1.1),
	}

	logger.Info("sending request to GigaChat")
	chatResp, err := r.client.Chat(ctx, chatReq)
	if err != nil {
		logger.WithError(err).Error("failed to get response from GigaChat")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	responseText := chatResp.Choices[0].Message.Content
	jsonData, err := r.extractJSON(responseText)
	if err != nil {
		logger.WithError(err).Error("failed to extract JSON from response")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	roadmapEntity, err := r.parseToRoadmap(jsonData)
	if err != nil {
		logger.WithError(err).Error("failed to parse JSON to roadmap")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	roadmapEntity.ID = primitive.NewObjectID()
	roadmapEntity.CreatedAt = time.Now()
	roadmapEntity.UpdatedAt = time.Now()

	logger.WithFields(map[string]interface{}{
		"nodes_count": len(roadmapEntity.Nodes),
		"edges_count": len(roadmapEntity.Edges),
	}).Info("successfully generated roadmap content")

	return roadmapEntity, nil
}

func (r *RoadmapGigaChatWebapi) loadRoadmapExample() (string, error) {
	example, err := roadmapPrompts.ReadFile("prompts/roadmap_example.json")
	if err != nil {
		return "", fmt.Errorf("failed to read roadmap example: %w", err)
	}
	return string(example), nil
}

func (r *RoadmapGigaChatWebapi) buildGenerationPrompt(req *dto.GenerateRoadmapDTO, example string) (string, error) {
	promptTmpl, err := roadmapPrompts.ReadFile("prompts/generation_prompt.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to read generation prompt template: %w", err)
	}

	tmpl, err := template.New("roadmapPrompt").Parse(string(promptTmpl))
	if err != nil {
		return "", fmt.Errorf("failed to parse generation prompt template: %w", err)
	}

	type PromptData struct {
		Topic       string
		Description string
		Content     string
		Complexity  string
		Example     string
	}

	data := PromptData{
		Topic:       req.Topic,
		Description: req.Description,
		Content:     req.Content,
		Complexity:  req.Complexity,
		Example:     example,
	}

	buf := &strings.Builder{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute generation prompt template: %w", err)
	}

	return buf.String(), nil
}

func (r *RoadmapGigaChatWebapi) extractJSON(text string) (string, error) {
	result := gjson.Parse(text)

	if result.Exists() && result.Type != gjson.Null {
		return result.String(), nil
	}

	return "", fmt.Errorf("valid JSON not found in response")
}

func (r *RoadmapGigaChatWebapi) parseToRoadmap(jsonData string) (*entities.Roadmap, error) {
	var roadmap entities.Roadmap
	if err := json.Unmarshal([]byte(jsonData), &roadmap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal roadmap: %w", err)
	}
	return &roadmap, nil
}

func float64Ptr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}
