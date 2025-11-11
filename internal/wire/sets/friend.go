package sets

import (
	"github.com/google/wire"

	friendHandlers "github.com/F0urward/proftwist-backend/services/friend/delivery/http"
	friendRepository "github.com/F0urward/proftwist-backend/services/friend/repository"
	friendUsecase "github.com/F0urward/proftwist-backend/services/friend/usecase"
)

var FriendSet = wire.NewSet(
	friendRepository.NewFriendRepository,
	friendUsecase.NewFriendUsecase,
	friendHandlers.NewFriendHandlers,
)
