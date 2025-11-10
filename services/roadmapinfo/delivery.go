package roadmapinfo

import (
	"net/http"
)

type Handlers interface {
	GetAllPublic(w http.ResponseWriter, r *http.Request)
	GetAllPublicByCategoryID(w http.ResponseWriter, r *http.Request)
	GetAllByUserID(w http.ResponseWriter, r *http.Request)
	GetByID(http.ResponseWriter, *http.Request)
	GetByRoadmapID(w http.ResponseWriter, r *http.Request)
	CreatePrivate(http.ResponseWriter, *http.Request)
	UpdatePrivate(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Fork(w http.ResponseWriter, r *http.Request)
	Publish(w http.ResponseWriter, r *http.Request)
	Subscribe(w http.ResponseWriter, r *http.Request)
	Unsubscribe(w http.ResponseWriter, r *http.Request)
	GetSubscribedRoadmaps(w http.ResponseWriter, r *http.Request)
}
