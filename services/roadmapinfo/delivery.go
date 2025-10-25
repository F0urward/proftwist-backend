package roadmapinfo

import (
	"net/http"
)

type Handlers interface {
	GetAll(http.ResponseWriter, *http.Request)
	GetAllByCategoryID(w http.ResponseWriter, r *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	GetByRoadmapID(w http.ResponseWriter, r *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}
