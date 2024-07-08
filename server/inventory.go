package server

import (
	"log"

	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryHandler"
	inventoryPb "github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryPb"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryUsecase"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/grpccon"
)

func (s *server) inventoryService() {
	repo := inventoryRepository.NewInventoryService(s.db)
	usecase := inventoryUsecase.NewInvertoryUsecase(repo)
	httpHandler := inventoryHandler.NewInvertoryHttpHandler(s.cfg, usecase)
	grpcHandler := inventoryHandler.NewInvertoryGrpcHandler(usecase)
	queueHandler := inventoryHandler.NewInvertoryQueueHandler(usecase)

	// gRPC
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)

		inventoryPb.RegisterInventoryGrpcServiceServer(grpcServer, grpcHandler)

		log.Printf("Inventory gRPC server listening on %s", s.cfg.Grpc.InventoryUrl)
		grpcServer.Serve(lis)
	}()

	_ = grpcHandler
	_ = queueHandler

	inventory := s.app.Group("/inventory_v1")

	// Health Check
	inventory.GET("", s.healthCheckService)
	inventory.GET("/inventory/:player_id", httpHandler.FindPlayerItems, s.middleware.JwtAuthorization, s.middleware.PlayerIdParamValidation)
}
