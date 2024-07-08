package inventoryHandler

import (
	"context"
	"net/http"

	"github.com/bonxatiwat/bonx-shop-tutorial/config"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryUsecase"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/request"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	InventoryHttpHandlerService interface {
		FindPlayerItems(c echo.Context) error
	}

	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInvertoryHttpHandler(cfg *config.Config, inventoryUsecase inventoryUsecase.InventoryUsecaseService) InventoryHttpHandlerService {
	return &inventoryHttpHandler{cfg, inventoryUsecase}
}

func (h *inventoryHttpHandler) FindPlayerItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	playerId := c.Param("player_id")

	req := new(inventory.InventorySearchReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.inventoryUsecase.FindPlayerItems(ctx, h.cfg.Paginate.InventoryNextPageBasedUrl, playerId, req)
	if err != nil {
		return response.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
