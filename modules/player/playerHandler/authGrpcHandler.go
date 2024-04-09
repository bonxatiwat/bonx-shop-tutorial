package playerHandler

import "github.com/bonxatiwat/bonx-shop-tutorial/modules/player/playerUsecase"

type (
	playerGrpcHandlerService struct {
		playerUsecase playerUsecase.PlayerUsecaseService
	}
)

func NewPlayerGrpcHandler(playerUsecase playerUsecase.PlayerUsecaseService) playerGrpcHandlerService {
	return playerGrpcHandlerService{playerUsecase}
}
