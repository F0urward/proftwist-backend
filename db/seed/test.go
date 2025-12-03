package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"path/filepath"

// 	"github.com/google/uuid"
// )

// type NodePosition struct {
// 	X float64 `json:"x"`
// 	Y float64 `json:"y"`
// }

// type NodeData struct {
// 	Label      string `json:"label"`
// 	Type       string `json:"type"`
// 	IsSelected bool   `json:"isSelected"`
// }

// type NodeMeasured struct {
// 	Width  float64 `json:"width"`
// 	Height float64 `json:"height"`
// }

// type RoadmapNode struct {
// 	ID          string       `json:"id"`
// 	Type        string       `json:"type"`
// 	Position    NodePosition `json:"position"`
// 	Data        NodeData     `json:"data"`
// 	Measured    NodeMeasured `json:"measured"`
// 	Selected    bool         `json:"selected"`
// 	Dragging    bool         `json:"dragging"`
// 	Description string       `json:"description"`
// 	Materials   []string     `json:"materials"`
// }

// type RoadmapEdge struct {
// 	ID     string                 `json:"id"`
// 	Source string                 `json:"source"`
// 	Target string                 `json:"target"`
// 	Type   string                 `json:"type"`
// 	Data   map[string]interface{} `json:"data"`
// }

// type Roadmap struct {
// 	Nodes []RoadmapNode `json:"nodes"`
// 	Edges []RoadmapEdge `json:"edges"`
// }

// type RoadmapJSON struct {
// 	Name string  `json:"name"`
// 	Data Roadmap `json:"data"`
// }

// type RoadmapCollection struct {
// 	Roadmaps []RoadmapJSON `json:"roadmaps"`
// }

// func RegenerateRoadmapIDs(rm *Roadmap) {
// 	idMap := make(map[string]string)

// 	// Replace node IDs
// 	for i, n := range rm.Nodes {
// 		newID := uuid.New().String()
// 		idMap[n.ID] = newID
// 		rm.Nodes[i].ID = newID
// 	}

// 	// Update edges with mapped IDs
// 	for i, e := range rm.Edges {
// 		if newSrc, ok := idMap[e.Source]; ok {
// 			rm.Edges[i].Source = newSrc
// 		}
// 		if newTgt, ok := idMap[e.Target]; ok {
// 			rm.Edges[i].Target = newTgt
// 		}
// 	}
// }

// func main() {
// 	inputPath := filepath.Join("data", "roadmaps.json")
// 	outputPath := filepath.Join("data", "roadmaps_corrected.json")

// 	// Read the original JSON
// 	raw, err := os.ReadFile(inputPath)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to read %s: %w", inputPath, err))
// 	}

// 	// Parse JSON into local structs
// 	var collection RoadmapCollection
// 	if err := json.Unmarshal(raw, &collection); err != nil {
// 		panic(fmt.Errorf("failed to unmarshal JSON: %w", err))
// 	}

// 	// Fix all maps
// 	for i := range collection.Roadmaps {
// 		RegenerateRoadmapIDs(&collection.Roadmaps[i].Data)
// 	}

// 	// Output JSON
// 	out, err := json.MarshalIndent(collection, "", "  ")
// 	if err != nil {
// 		panic(fmt.Errorf("failed to marshal JSON: %w", err))
// 	}

// 	if err := os.WriteFile(outputPath, out, 0644); err != nil {
// 		panic(fmt.Errorf("failed to write file: %w", err))
// 	}

// 	fmt.Println("✔ Corrected JSON saved to", outputPath)
// }
