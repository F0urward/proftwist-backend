package roadmap

import "net/http"

type Handlers interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetByAuthorID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	SearchByTitle(w http.ResponseWriter, r *http.Request)
	UpdatePrivacy(http.ResponseWriter, *http.Request)
}
