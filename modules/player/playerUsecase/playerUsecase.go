package playerUsecase

import (
	"context"
	"errors"

	"github.com/bonxatiwat/bonx-shop-tutorial/modules/player"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/player/playerRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error)
	}

	playerUsecase struct {
		playerRepository playerRepository.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepository playerRepository.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{playerRepository}
}

func (u *playerUsecase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error) {
	if !u.playerRepository.IsUniquePlayer(pctx, req.Email, req.Username) {
		return "", errors.New("error: email or username already exist")
	}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error: failed to hash password")
	}

	// Insert one player
	playerId, err := u.playerRepository.InsertOnePlayer(pctx, &player.Player{
		Email:     req.Email,
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		PlayerRoles: []player.PlayerRole{
			{
				RoleTitle: "player",
				RoleCode:  0,
			},
		},
	})
	return playerId.Hex(), nil
}
