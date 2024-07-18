package paymentUsecase

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/bonxatiwat/bonx-shop-tutorial/config"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/item"
	itemPb "github.com/bonxatiwat/bonx-shop-tutorial/modules/item/itemPb"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/payment"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/payment/paymentRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/player"
	"github.com/bonxatiwat/bonx-shop-tutorial/pkg/queue"
)

type (
	PaymentUsecaseService interface {
		FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error
		GetOffset(pctx context.Context) (int64, error)
		UpsertOffset(pctx context.Context, offset int64) error
		BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error)
		BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes)
	}

	paymentUsecase struct {
		paymentRepository paymentRepository.PaymentRepositoryService
	}
)

func NewPaymentUsecase(paymentRepository paymentRepository.PaymentRepositoryService) PaymentUsecaseService {
	return &paymentUsecase{paymentRepository: paymentRepository}
}

func (u *paymentUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.paymentRepository.GetOffset(pctx)
}

func (u *paymentUsecase) UpsertOffset(pctx context.Context, offset int64) error {
	return u.paymentRepository.UpsertOffset(pctx, offset)
}

func (u *paymentUsecase) PaymentConsumer(pctx context.Context, cfg *config.Config) (sarama.PartitionConsumer, error) {
	worker, err := queue.ConnectConsumer([]string{cfg.Kafka.Url}, cfg.Kafka.ApiKey, cfg.Kafka.Secret)
	if err != nil {
		return nil, err
	}

	offset, err := u.paymentRepository.GetOffset(pctx)
	if err != nil {
		return nil, err
	}

	consumer, err := worker.ConsumePartition("payment", 0, offset)
	if err != nil {
		log.Println("Trying to set offset as 0")
		consumer, err = worker.ConsumePartition("payment", 0, 0)
		if err != nil {
			log.Println("Error: PaymentConsumer failed: ", err.Error())
			return nil, err
		}
	}

	return consumer, nil
}

func (u *paymentUsecase) BuyOrSellConsumer(pctx context.Context, key string, cfg *config.Config, resCh chan<- *payment.PaymentTransferRes) {
	consumer, err := u.PaymentConsumer(pctx, cfg)
	if err != nil {
		resCh <- nil
		return
	}
	defer consumer.Close()

	log.Println("Start BuyOrSellConsumer ...")

	select {
	case err := <-consumer.Errors():
		log.Println("Error: BuyOrSellConsumer failed: ", err.Error())
		resCh <- nil
		return
	case msg := <-consumer.Messages():
		if string(msg.Key) == key {
			u.UpsertOffset(pctx, msg.Offset+1)

			req := new(payment.PaymentTransferRes)

			if err := queue.DecodeMessage(req, msg.Value); err != nil {
				resCh <- nil
				return
			}

			resCh <- req
			log.Printf("BuyOrSellConsumer | Topic(%s)| Offset(%d) Message(%s) \n", msg.Topic, msg.Offset, string(msg.Value))
		}
	}
}

func (u *paymentUsecase) BuyItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}

	stage1 := make([]*payment.PaymentTransferRes, 0)
	for _, item := range req.Items {
		u.paymentRepository.DockedPlayerMoney(pctx, cfg, &player.CreatePlayerTransactionReq{
			PlayerId: playerId,
			Amount:   -item.Price,
		})

		resCh := make(chan *payment.PaymentTransferRes)

		go u.BuyOrSellConsumer(pctx, "buy", cfg, resCh)

		res := <-resCh
		if res != nil {
			stage1 = append(stage1, res)
		}
	}

	for _, s1 := range stage1 {
		if s1.Error != "" {
			for _, ss1 := range stage1 {
				u.paymentRepository.RollbackTransaction(pctx, cfg, &player.RollbackPlayerTransactionReq{
					TransactionId: ss1.TransactionId,
				})
			}
			return nil, errors.New("error: buy item failed")
		}
	}

	return stage1, nil
}

func (u *paymentUsecase) SellItem(pctx context.Context, cfg *config.Config, playerId string, req *payment.ItemServiceReq) ([]*payment.PaymentTransferRes, error) {
	if err := u.FindItemsInIds(pctx, cfg.Grpc.ItemUrl, req.Items); err != nil {
		return nil, err
	}
	return nil, nil
}

func (u *paymentUsecase) FindItemsInIds(pctx context.Context, grpcUrl string, req []*payment.ItemServiceReqDatum) error {

	setIds := make(map[string]bool)
	for _, v := range req {
		if !setIds[v.ItemId] {
			setIds[v.ItemId] = true
		}
	}

	itemData, err := u.paymentRepository.FindItemsInIds(pctx, grpcUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for k := range setIds {
				itemIds = append(itemIds, k)
			}
			return itemIds
		}(),
	})
	if err != nil {
		log.Printf("Error: FindItemsInIds failed: %s", err.Error())
		return errors.New("error: item not found")
	}

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

	for i := range req {
		if _, ok := itemMaps[req[i].ItemId]; !ok {
			log.Printf("Error: FindItemsInIds failed: %s", err.Error())
			return errors.New("error: items not found")
		}
		req[i].Price = itemMaps[req[i].ItemId].Price
	}

	return nil
}
