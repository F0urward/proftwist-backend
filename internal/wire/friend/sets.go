package friend

import (
	"github.com/google/wire"

	friendHandlers "github.com/F0urward/proftwist-backend/services/friend/delivery/http"
	friendRepository "github.com/F0urward/proftwist-backend/services/friend/repository"
	friendUsecase "github.com/F0urward/proftwist-backend/services/friend/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	chatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
)

var FriendSet = wire.NewSet(
	friendRepository.NewFriendRepository,
	friendUsecase.NewFriendUsecase,
	friendHandlers.NewFriendHandlers,
	friendHandlers.NewFriendHttpRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	chatClient.NewChatClient,
	authClient.NewAuthClient,
)
