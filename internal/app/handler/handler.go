package handler

import (
	"github.com/harunnryd/tempolalu/config"
	"github.com/harunnryd/tempolalu/internal/app/handler/transaction"
	"github.com/harunnryd/tempolalu/internal/app/usecase"
)

type Handler interface {
	Transaction() transaction.Transaction
}

type handler struct {
	transaction transaction.Transaction
}

func NewHandler(cfg config.Config, usecase usecase.Usecase) Handler {
	h := new(handler)

	h.transaction = transaction.New(usecase)

	return h
}

func (h *handler) Transaction() transaction.Transaction {
	return h.transaction
}
