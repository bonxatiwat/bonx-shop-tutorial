package item

import "github.com/bonxatiwat/bonx-shop-tutorial/modules/models"

type (
	CreateItemReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"price"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
		Damage   int     `json:"damage" validate:"required,max=255"`
	}

	ItemShowCase struct {
		ItemId   string `json:"item_id"`
		Title    string `json:"title"`
		Price    string `json:"price"`
		Damage   int    `json:"damage"`
		ImageUrl string `json:"image_url"`
	}

	ItemSearchReq struct {
		Title string `json:"title"`
		models.PaginateReq
	}

	ItemUpdateReq struct {
		Title    string  `json:"title" validate:"required,max=64"`
		Price    float64 `json:"price" validate:"required"`
		ImageUrl string  `json:"image_url" validate:"required,max=255"`
		Damage   int     `json:"damage" validate:"required,max=255"`
	}

	EnableOrDisableItemReq struct {
		UsageStauts bool `json:"staus"`
	}
)
