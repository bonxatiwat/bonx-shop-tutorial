package server

import (
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryHandler"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryUsecase"
)

func (s *server) inventoryService() {
	repo := inventoryRepository.NewInventoryService(s.db)
	usecase := inventoryUsecase.NewInvertoryUsecase(repo)
	httpHandler := inventoryHandler.NewInvertoryHttpHandler(s.cfg, usecase)
	grpcHandler := inventoryHandler.NewInvertoryGrpcHandler(usecase)
	queueHandler := inventoryHandler.NewInvertoryQueueHandler(usecase)

	_ = httpHandler
	_ = grpcHandler
	_ = queueHandler

	inventory := s.app.Group("/inventory_v1")

	// Health Check
	_ = inventory
}
