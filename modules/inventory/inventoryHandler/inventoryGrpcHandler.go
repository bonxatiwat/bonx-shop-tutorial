package inventoryHandler

import (
	"context"

	inventoryPb "github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryPb"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryUsecase"
)

type (
	inventoryGrpcHandler struct {
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
		inventoryPb.UnsafeInventoryGrpcServiceServer
	}
)

func NewInvertoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{inventoryUsecase: inventoryUsecase}
}

func (g *inventoryGrpcHandler) IsAvaliableToSell(ctx context.Context, req *inventoryPb.IsAvaliableToSellReq) (*inventoryPb.IsAvaliableToSellRes, error) {
	return nil, nil
}
