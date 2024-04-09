package inventoryHandler

import "github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryUsecase"

type (
	InventoryQueueHandlerService interface{}

	inventoryQueueHandler struct {
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInvertoryQueueHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryQueueHandlerService {
	return &inventoryQueueHandler{inventoryUsecase}
}
