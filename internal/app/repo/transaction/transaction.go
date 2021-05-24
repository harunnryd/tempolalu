package transaction

import (
	"context"

	"github.com/harunnryd/tempolalu/config"
	"github.com/harunnryd/tempolalu/internal/app/model"
	"gorm.io/gorm"
)

type Transaction interface {
	Create(context.Context, int, int) (model.Transaction, error)
	GetByID(context.Context, int) (model.Transaction, error)
	DeleteByID(context.Context, int) error
}

type transaction struct {
	cfg      config.Config
	ormPgSQL *gorm.DB
}

func New(cfg config.Config, ormPgSQL *gorm.DB) Transaction {
	return &transaction{cfg: cfg, ormPgSQL: ormPgSQL}
}

func (t *transaction) Create(ctx context.Context, productID int, quantity int) (transactionModel model.Transaction, err error) {
	transactionModel = model.Transaction{ProductID: productID, Quantity: quantity}
	err = t.ormPgSQL.WithContext(ctx).Create(&transactionModel).Error

	return
}

func (t *transaction) GetByID(ctx context.Context, id int) (transactionModel model.Transaction, err error) {
	err = t.ormPgSQL.WithContext(ctx).Where("id = ?", id).Find(&transactionModel).Error

	return
}

func (t *transaction) DeleteByID(ctx context.Context, transactionID int) (err error) {
	err = t.ormPgSQL.WithContext(ctx).Exec(`DELETE FROM transactions WHERE id = ?`, transactionID).Error
	return
}
