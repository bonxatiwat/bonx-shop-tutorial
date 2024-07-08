package inventoryUsecase

import (
	"context"
	"fmt"

	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/item"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/models"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryUsecaseService interface {
		FindPlayerItems(pctx context.Context, basePaginateUrl, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
	}

	inventoryUsecase struct {
		inventoryRepository inventoryRepository.InventoryRepositoryService
	}
)

func NewInvertoryUsecase(inventoryRepository inventoryRepository.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{inventoryRepository}
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, basePaginateUrl, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
	// Filter
	filter := bson.D{}

	if req.Start != "" {
		filter = append(filter, bson.E{"_id", bson.D{{"$gt", utils.ConvertToObjectId(req.Start)}}})
	}
	filter = append(filter, bson.E{"player_id", playerId})

	// Option
	opts := make([]*options.FindOptions, 0)

	opts = append(opts, options.Find().SetSort(bson.D{{"_id", 1}}))
	opts = append(opts, options.Find().SetLimit(int64(req.Limit)))

	// Find
	inventoryData, err := u.inventoryRepository.FindPlayerItems(pctx, filter, opts)
	if err != nil {
		return nil, err
	}

	results := make([]*inventory.ItemInInventory, 0)
	for _, v := range inventoryData {
		results = append(results, &inventory.ItemInInventory{
			InventoryId: v.Id,
			PlayerId:    v.PlayerId,
			ItemShowCase: &item.ItemShowCase{
				ItemId: v.Id,
			},
		})
	}

	// Count
	total, err := u.inventoryRepository.CountPlayerItems(pctx, playerId)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &models.PaginateRes{
			Data:  make([]*inventory.ItemInInventory, 0),
			Total: 0,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s/%s?limit=%d", basePaginateUrl, playerId, req.Limit),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", basePaginateUrl, playerId, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryId,
			Href:  fmt.Sprintf("%s/%s?limit=%d&start=%s", basePaginateUrl, playerId, req.Limit, results[len(results)-1].ItemId),
		},
	}, nil
}
