package server

import (
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/payment/paymentHandler"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/payment/paymentRepository"
	"github.com/bonxatiwat/bonx-shop-tutorial/modules/payment/paymentUsecase"
)

func (s *server) paymentService() {
	repo := paymentRepository.NewPaymentRepository(s.db)
	usecase := paymentUsecase.NewPaymentUsecase(repo)
	httpHandler := paymentHandler.NewPaymentHttpHandler(s.cfg, usecase)
	queueHandler := paymentHandler.NewPaymentQueueHandler(s.cfg, usecase)

	_ = httpHandler
	_ = queueHandler

	payment := s.app.Group("/payment_v1")

	// Health Check
	_ = payment
}
