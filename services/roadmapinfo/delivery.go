package roadmapinfo

import (
	"net/http"
)

type Handlers interface {
	GetAll(http.ResponseWriter, *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}
