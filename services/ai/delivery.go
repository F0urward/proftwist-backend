package ai

import "net/http"

type Handlers interface {
	GenerateRoadmapNodeDescription(w http.ResponseWriter, r *http.Request)
	GenerateRoadmap(w http.ResponseWriter, r *http.Request)
}
