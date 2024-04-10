package authHandler

import (
	"context"

	authPb "github.com/bonxatiwat/bonx-shop-tutorial/modules/auth/authPb"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/auth/authUsecase"
)

type (
	authGrpcHandler struct {
		authPb.UnimplementedAuthGrpcServiceServer
		authUsecase authUsecase.AuthUsecaseService
	}
)

func NewAuthGrpcHandler(authUsecase authUsecase.AuthUsecaseService) *authGrpcHandler {
	return &authGrpcHandler{
		authUsecase: authUsecase,
	}
}

func (g *authGrpcHandler) CrednetialSearch(ctx context.Context, req *authPb.CredentialSearchReq) (*authPb.CredentialSearchRes, error) {
	return nil, nil
}

func (g *authGrpcHandler) RolesCount(ctx context.Context, req *authPb.RolesCountReq) (*authPb.RolesCountRes, error) {
	return nil, nil
}
