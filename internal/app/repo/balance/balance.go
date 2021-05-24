package balance

import (
	"context"
	"database/sql"

	"github.com/harunnryd/tempolalu/config"
	"github.com/harunnryd/tempolalu/internal/app/model"
	"gorm.io/gorm"
)

type Balance interface {
	Debit(context.Context, int, float64) (model.Balance, error)
	Credit(context.Context, int, float64) (model.Balance, error)
}

type balance struct {
	cfg      config.Config
	ormPgSQL *gorm.DB
}

const DebitQuery = `UPDATE balances SET nominal = nominal - @nominal WHERE user_id = @user_id RETURNING id, user_id, nominal, created_at, updated_at, deleted_at`
const CreditQuery = `UPDATE balances SET nominal = nominal + @nominal WHERE user_id = @user_id RETURNING id, user_id, nominal, created_at, updated_at, deleted_at`

func New(cfg config.Config, ormPgSQL *gorm.DB) Balance {
	return &balance{cfg: cfg, ormPgSQL: ormPgSQL}
}

func (b *balance) Debit(ctx context.Context, userID int, nominal float64) (balanceModel model.Balance, err error) {
	rows, err := b.ormPgSQL.WithContext(ctx).Raw(DebitQuery, sql.Named("nominal", nominal), sql.Named("user_id", userID)).Rows()

	for rows.Next() {
		b.ormPgSQL.ScanRows(rows, &balanceModel)
	}
	return
}

func (b *balance) Credit(ctx context.Context, userID int, nominal float64) (balanceModel model.Balance, err error) {
	rows, err := b.ormPgSQL.WithContext(ctx).Raw(CreditQuery, sql.Named("nominal", nominal), sql.Named("user_id", userID)).Rows()

	for rows.Next() {
		b.ormPgSQL.ScanRows(rows, &balanceModel)
	}
	return
}
