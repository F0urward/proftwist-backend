package material

import (
	"net/http"
)

type Handlers interface {
	CreateMaterial(w http.ResponseWriter, r *http.Request)
	DeleteMaterial(w http.ResponseWriter, r *http.Request)
	GetMaterialsByNode(w http.ResponseWriter, r *http.Request)
	GetUserMaterials(w http.ResponseWriter, r *http.Request)
}
