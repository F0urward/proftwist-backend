package roadmap

import "net/http"

type Handlers interface {
	GetByIDWithProgress(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Generate(w http.ResponseWriter, r *http.Request)
	CreateMaterial(w http.ResponseWriter, r *http.Request)
	DeleteMaterial(w http.ResponseWriter, r *http.Request)
	GetMaterialsByNode(w http.ResponseWriter, r *http.Request)
	UpdateNodeProgress(w http.ResponseWriter, r *http.Request)
}
