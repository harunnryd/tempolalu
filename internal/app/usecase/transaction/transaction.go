package transaction

import (
	"context"

	"github.com/harunnryd/tempolalu/internal/app/model"
	"github.com/harunnryd/tempolalu/internal/app/repo"
)

type Transaction interface {
	CreateTransaction(context.Context, int, int, int, int, float64) (model.Transaction, error)
	DeleteTransactionByID(context.Context, int) error
}

type transaction struct {
	repo repo.Repo
}

func New(repo repo.Repo) Transaction {
	return &transaction{repo: repo}
}

func (t *transaction) DeleteTransactionByID(ctx context.Context, id int) (err error) {
	err = t.repo.Transaction().DeleteByID(ctx, id)
	if err != nil {
		return
	}

	mutationDebitModel, err := t.repo.Mutation().GetMutationDebitByTransactionID(ctx, id)
	if err != nil {
		return
	}

	_, err = t.repo.Balance().Credit(ctx, mutationDebitModel.OriginID, mutationDebitModel.Amount)
	if err != nil {
		return
	}

	_, err = t.repo.Balance().Debit(ctx, mutationDebitModel.DestinationID, mutationDebitModel.Amount)
	if err != nil {
		return
	}

	err = t.repo.Mutation().DeleteByTransactionID(ctx, id)
	if err != nil {
		return
	}

	err = t.repo.MutationTransaction().DeleteByTransactionID(ctx, id)

	return
}

const (
	PAID          = 4
	POCKET_CREDIT = 1
	POCKET_DEBIT  = 2
	CHANNEL_TF    = 1
)

func (t *transaction) CreateTransaction(ctx context.Context, originID int, destinationID int, productID int, quantity int, amount float64) (transactionModel model.Transaction, err error) {
	transactionModel, err = t.repo.Transaction().Create(ctx, productID, quantity)
	if err != nil {
		return
	}

	err = t.credit(ctx, originID, destinationID, transactionModel.ID, quantity, amount)
	if err != nil {
		return
	}

	err = t.debit(ctx, originID, destinationID, transactionModel.ID, quantity, amount)

	return
}

func (t *transaction) debit(ctx context.Context, originID int, destinationID int, transactionID int, quantity int, amount float64) (err error) {
	_, err = t.repo.Balance().Debit(ctx, originID, amount)
	if err != nil {
		return
	}

	mutationModel, err := t.repo.Mutation().Create(ctx, originID, destinationID, amount, PAID, POCKET_DEBIT, CHANNEL_TF)
	if err != nil {
		return
	}

	_, err = t.repo.MutationTransaction().Create(ctx, mutationModel.ID, transactionID)

	return
}

func (t *transaction) credit(ctx context.Context, originID int, destinationID int, transactionID int, quantity int, amount float64) (err error) {
	_, err = t.repo.Balance().Credit(ctx, destinationID, amount)
	if err != nil {
		return
	}

	mutationModel, err := t.repo.Mutation().Create(ctx, destinationID, originID, amount, PAID, POCKET_CREDIT, CHANNEL_TF)
	if err != nil {
		return
	}

	_, err = t.repo.MutationTransaction().Create(ctx, mutationModel.ID, transactionID)

	return
}
