package mutation

import (
	"context"
	"database/sql"

	"github.com/harunnryd/tempolalu/config"
	"github.com/harunnryd/tempolalu/internal/app/model"
	"gorm.io/gorm"
)

type Mutation interface {
	Create(context.Context, int, int, float64, int, int, int) (model.Mutation, error)
	GetMutationDebitByTransactionID(context.Context, int) (model.Mutation, error)
	UpdateStatusByTransactionID(context.Context, int, int) error
	DeleteByTransactionID(context.Context, int) error
}

type mutation struct {
	cfg      config.Config
	ormPgSQL *gorm.DB
}

func New(cfg config.Config, ormPgSQL *gorm.DB) Mutation {
	return &mutation{cfg: cfg, ormPgSQL: ormPgSQL}
}

func (m *mutation) Create(ctx context.Context, originID int, destinationID int, amount float64, status int, pocket int, channel int) (mutationModel model.Mutation, err error) {
	mutationModel = model.Mutation{
		OriginID:      originID,
		DestinationID: destinationID,
		Amount:        amount,
		Status:        status,
		Pocket:        pocket,
		Channel:       channel,
	}

	if err = m.ormPgSQL.WithContext(ctx).Create(&mutationModel).Error; err != nil {
		return
	}

	return
}

const UpdateStatusByTransactionIDQuery = `UPDATE mutations SET status = @status WHERE id IN (SELECT id FROM transactions WHERE id = @transaction_id)`

func (m *mutation) UpdateStatusByTransactionID(ctx context.Context, id int, status int) (err error) {
	err = m.ormPgSQL.WithContext(ctx).Raw(UpdateStatusByTransactionIDQuery).Error
	if err != nil {
		return
	}

	return
}

const DeleteByTransactionIDQuery = `DELETE FROM mutations WHERE id IN (SELECT mutation_id FROM mutation_transaction WHERE transaction_id = @transaction_id)`

func (m *mutation) DeleteByTransactionID(ctx context.Context, transactionID int) (err error) {
	err = m.ormPgSQL.WithContext(ctx).Exec(DeleteByTransactionIDQuery, sql.Named("transaction_id", transactionID)).Error
	return
}

const GetMutationDebitByTransactionIDQuery = `SELECT * FROM mutations WHERE id IN (SELECT mutation_id FROM mutation_transaction WHERE transaction_id = @transaction_id) AND pocket = 2 LIMIT 1`

func (m *mutation) GetMutationDebitByTransactionID(ctx context.Context, transactionID int) (mutationModel model.Mutation, err error) {
	rows, err := m.ormPgSQL.WithContext(ctx).Raw(GetMutationDebitByTransactionIDQuery).Rows()

	for rows.Next() {
		m.ormPgSQL.ScanRows(rows, &mutationModel)
	}

	return
}
