package usecase

import (
	"github.com/harunnryd/tempolalu/internal/app/repo"
	"github.com/harunnryd/tempolalu/internal/app/usecase/transaction"
)

type Usecase interface {
	Transaction() transaction.Transaction
}

type usecase struct {
	transaction transaction.Transaction
}

func NewUsecase(repo repo.Repo) Usecase {
	u := new(usecase)

	u.transaction = transaction.New(repo)

	return u
}

func (u *usecase) Transaction() transaction.Transaction {
	return u.transaction
}
