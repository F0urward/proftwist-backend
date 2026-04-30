package dto

type GenerateRoadmapNodeDescriptionRequestDTO struct {
	RoadmapID          string `json:"roadmap_id,omitempty"`
	NodeID             string `json:"node_id"`
	NodeLabel          string `json:"node_label"`
	NodeType           string `json:"node_type,omitempty"`
	CurrentDescription string `json:"current_description,omitempty"`
}

type GenerateRoadmapNodeDescriptionResponseDTO struct {
	Description string `json:"description"`
}
