package dto

type GenerateRoadmapNodeDescriptionRequestDTO struct {
	RoadmapID          string   `json:"roadmap_id,omitempty"`
	NodeID             string   `json:"node_id"`
	NodeLabel          string   `json:"node_label"`
	NodeType           string   `json:"node_type,omitempty"`
	CurrentDescription string   `json:"current_description,omitempty"`
	Provider           string   `json:"provider,omitempty"`
	Model              string   `json:"model,omitempty"`

	RoadmapName      string   `json:"-"`
	RootNodeLabel    string   `json:"-"`
	RootNodeType     string   `json:"-"`
	SiblingLabels    []string `json:"-"`
	ChildLabels      []string `json:"-"`
	TotalNodeCount   int      `json:"-"`
}

type GenerateRoadmapRequestDTO struct {
	RoadmapID string `json:"roadmap_id,omitempty"`
	Prompt    string `json:"prompt"`
	Provider  string `json:"provider,omitempty"`
	Model     string `json:"model,omitempty"`
}

type GenerateRoadmapResponseDTO struct {
	Roadmap RoadmapGraphDTO `json:"roadmap"`
}

type RoadmapGraphDTO struct {
	Nodes []RoadmapNodeDTO `json:"nodes"`
	Edges []RoadmapEdgeDTO `json:"edges"`
}

type RoadmapNodeDTO struct {
	ID          string             `json:"id"`
	Type        string             `json:"type"`
	Position    PositionDTO        `json:"position"`
	Data        RoadmapNodeDataDTO `json:"data"`
	Description string             `json:"description,omitempty"`
}

type PositionDTO struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type RoadmapNodeDataDTO struct {
	Label string `json:"label"`
	Type  string `json:"type"`
}

type RoadmapEdgeDTO struct {
	ID     string             `json:"id"`
	Source string             `json:"source"`
	Target string             `json:"target"`
	Type   string             `json:"type,omitempty"`
	Data   RoadmapEdgeDataDTO `json:"data,omitempty"`
}

type RoadmapEdgeDataDTO struct {
	Variant string `json:"variant,omitempty"`
}
