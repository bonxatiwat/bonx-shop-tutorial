package inventoryUsecase

import (
	"context"
	"fmt"

	"github.com/bonxatiwat/bonx-shop-tutorial/config"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/inventory/inventoryRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/item"
	itemPb "github.com/bonxatiwat/bonx-shop-tutorial/modules/item/itemPb"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/models"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/payment"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	InventoryUsecaseService interface {
		FindPlayerItems(pctx context.Context, cfg *config.Config, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RemovePlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq)
		RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq)
	}

	inventoryUsecase struct {
		inventoryRepository inventoryRepository.InventoryRepositoryService
	}
)

func NewInventoryUsecase(inventoryRepository inventoryRepository.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{inventoryRepository: inventoryRepository}
}

func (u *inventoryUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.inventoryRepository.GetOffset(pctx)
}

func (u *inventoryUsecase) UpsertOffset(pctx context.Context, offset int64) error {
	return u.inventoryRepository.UpsertOffset(pctx, offset)
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, cfg *config.Config, playerId string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
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

	if len(inventoryData) == 0 {
		return &models.PaginateRes{
			Data:  make([]*inventory.ItemInInventory, 0),
			Total: 0,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.ItemNextPageBasedUrl, playerId, req.Limit),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	itemData, err := u.inventoryRepository.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for _, v := range inventoryData {
				itemIds = append(itemIds, v.ItemId)
			}
			return itemIds
		}(),
	})

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemId:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			ImageUrl: v.ImageUrl,
			Damage:   int(v.Damage),
		}
	}

	results := make([]*inventory.ItemInInventory, 0)
	for _, v := range inventoryData {
		results = append(results, &inventory.ItemInInventory{
			InventoryId: v.Id,
			PlayerId:    v.PlayerId,
			ItemShowCase: &item.ItemShowCase{
				ItemId:   v.ItemId,
				Title:    itemMaps[v.ItemId].Title,
				Price:    itemMaps[v.ItemId].Price,
				Damage:   itemMaps[v.ItemId].Damage,
				ImageUrl: itemMaps[v.ItemId].ImageUrl,
			},
		})
	}

	// Count
	total, err := u.inventoryRepository.CountPlayerItems(pctx, playerId)
	if err != nil {
		return nil, err
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.ItemNextPageBasedUrl, playerId, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryId,
			Href:  fmt.Sprintf("%s/%s?limit=%d&start=%s", cfg.Paginate.ItemNextPageBasedUrl, playerId, req.Limit, results[len(results)-1].ItemId),
		},
	}, nil
}

func (u *inventoryUsecase) AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {
	inventoryId, err := u.inventoryRepository.InsertOnePlayerItem(pctx, &inventory.Inventory{
		PlayerId: req.PlayerId,
		ItemId:   req.ItemId,
	})
	if err != nil {
		u.inventoryRepository.AddPlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryId:   inventoryId.Hex(),
			TransactionId: "",
			PlayerId:      req.PlayerId,
			ItemId:        req.ItemId,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}
	u.inventoryRepository.AddPlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryId:   inventoryId.Hex(),
		TransactionId: "",
		PlayerId:      req.PlayerId,
		ItemId:        req.ItemId,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RemovePlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {

}

func (u *inventoryUsecase) RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) {
	u.inventoryRepository.DeleteOneInventory(pctx, req.InventoryId)
}

func (u *inventoryUsecase) RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventoryReq) {
	u.inventoryRepository.InsertOnePlayerItem(pctx, &inventory.Inventory{
		PlayerId: req.PlayerId,
		ItemId:   req.ItemId,
	})
}
