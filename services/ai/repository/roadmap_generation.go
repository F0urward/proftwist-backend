package repository

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/F0urward/proftwist-backend/services/ai/dto"
)

type roadmapGenerationEnvelope struct {
	Roadmap          *dto.RoadmapGraphDTO `json:"roadmap,omitempty"`
	Graph            *dto.RoadmapGraphDTO `json:"graph,omitempty"`
	GeneratedRoadmap *dto.RoadmapGraphDTO `json:"generated_roadmap,omitempty"`
	Nodes            []dto.RoadmapNodeDTO `json:"nodes,omitempty"`
	Edges            []dto.RoadmapEdgeDTO `json:"edges,omitempty"`
}

func parseRoadmapGraphResponse(text string) (*dto.RoadmapGraphDTO, error) {
	jsonText, err := extractJSONObject(text)
	if err != nil {
		return nil, err
	}

	var envelope roadmapGenerationEnvelope
	if err := json.Unmarshal([]byte(jsonText), &envelope); err != nil {
		return nil, fmt.Errorf("failed to decode roadmap JSON: %w", err)
	}

	switch {
	case envelope.Roadmap != nil:
		return ensureRoadmapArrays(envelope.Roadmap), nil
	case envelope.Graph != nil:
		return ensureRoadmapArrays(envelope.Graph), nil
	case envelope.GeneratedRoadmap != nil:
		return ensureRoadmapArrays(envelope.GeneratedRoadmap), nil
	case envelope.Nodes != nil || envelope.Edges != nil:
		return ensureRoadmapArrays(&dto.RoadmapGraphDTO{
			Nodes: envelope.Nodes,
			Edges: envelope.Edges,
		}), nil
	default:
		return nil, fmt.Errorf("roadmap JSON does not contain nodes and edges")
	}
}

func ensureRoadmapArrays(graph *dto.RoadmapGraphDTO) *dto.RoadmapGraphDTO {
	if graph == nil {
		return &dto.RoadmapGraphDTO{
			Nodes: []dto.RoadmapNodeDTO{},
			Edges: []dto.RoadmapEdgeDTO{},
		}
	}
	if graph.Nodes == nil {
		graph.Nodes = []dto.RoadmapNodeDTO{}
	}
	if graph.Edges == nil {
		graph.Edges = []dto.RoadmapEdgeDTO{}
	}
	return graph
}

func extractJSONObject(text string) (string, error) {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```JSON")
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	if json.Valid([]byte(text)) {
		return text, nil
	}

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start < 0 || end < start {
		return "", fmt.Errorf("valid JSON object not found in provider response")
	}

	candidate := strings.TrimSpace(text[start : end+1])
	if !json.Valid([]byte(candidate)) {
		return "", fmt.Errorf("extracted provider response is not valid JSON")
	}

	return candidate, nil
}