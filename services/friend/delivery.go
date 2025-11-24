package friend

import "net/http"

type Handlers interface {
	GetFriends(w http.ResponseWriter, r *http.Request)
	DeleteFriend(w http.ResponseWriter, r *http.Request)
	GetFriendRequests(w http.ResponseWriter, r *http.Request)
	AcceptFriendRequest(w http.ResponseWriter, r *http.Request)
	RejectFriendRequest(w http.ResponseWriter, r *http.Request)
	CreateFriendRequest(w http.ResponseWriter, r *http.Request)
	DeleteFriendRequest(w http.ResponseWriter, r *http.Request)
	GetFriendshipStatus(w http.ResponseWriter, r *http.Request)
}
