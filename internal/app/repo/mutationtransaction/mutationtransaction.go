package mutationtransaction

import (
	"context"
	"database/sql"

	"github.com/harunnryd/tempolalu/config"
	"github.com/harunnryd/tempolalu/internal/app/model"
	"gorm.io/gorm"
)

type MutationTransaction interface {
	Create(context.Context, int, int) (model.MutationTransaction, error)
	DeleteByTransactionID(context.Context, int) error
}

type mutationTransaction struct {
	cfg      config.Config
	ormPgSQL *gorm.DB
}

func New(cfg config.Config, ormPgSQL *gorm.DB) MutationTransaction {
	return &mutationTransaction{cfg: cfg, ormPgSQL: ormPgSQL}
}

func (m *mutationTransaction) Create(ctx context.Context, mutationID int, transactionID int) (mtransactionModel model.MutationTransaction, err error) {
	mtransactionModel = model.MutationTransaction{MutationID: mutationID, TransactionID: transactionID}
	err = m.ormPgSQL.WithContext(ctx).Create(&mtransactionModel).Error

	return
}

const DeleteByTransactionIDQuery = `DELETE FROM mutation_transaction WHERE transaction_id = @transaction_id`

func (m *mutationTransaction) DeleteByTransactionID(ctx context.Context, transactionID int) (err error) {
	m.ormPgSQL.WithContext(ctx).Exec(DeleteByTransactionIDQuery, sql.Named("transaction_id", transactionID))

	return
}
