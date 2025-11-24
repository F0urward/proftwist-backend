package friend

import (
	"github.com/google/wire"

	friendGrpc "github.com/F0urward/proftwist-backend/services/friend/delivery/grpc"
	friendHttp "github.com/F0urward/proftwist-backend/services/friend/delivery/http"
	friendRepository "github.com/F0urward/proftwist-backend/services/friend/repository"
	friendUsecase "github.com/F0urward/proftwist-backend/services/friend/usecase"

	authClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/authclient"
	chatClient "github.com/F0urward/proftwist-backend/internal/infrastructure/client/chatclient"
	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
)

var FriendSet = wire.NewSet(
	friendRepository.NewFriendRepository,
	friendUsecase.NewFriendUsecase,
	friendHttp.NewFriendHandlers,
	friendHttp.NewFriendHttpRegistrar,
	friendGrpc.NewFriendServer,
	friendGrpc.NewFriendGrpcRegistrar,
)

var ClientsSet = wire.NewSet(
	db.NewDatabase,
	chatClient.NewChatClient,
	authClient.NewAuthClient,
)
