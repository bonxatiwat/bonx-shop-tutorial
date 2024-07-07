package inventory

import (
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/item"
)

type (
	UpdateInventoryReq struct {
		PlayerId string `json:"player_id" validate:"required,max=64"`
		ItemId   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInventory struct {
		IventoryId string `json:"inventory_id"`
		PlayerId   string `json:"player_id"`
		*item.ItemShowCase
	}
)
