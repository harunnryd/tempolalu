package repo

import (
	"log"

	"github.com/harunnryd/tempolalu/config"
	"github.com/harunnryd/tempolalu/internal/app/driver/db"
	"github.com/harunnryd/tempolalu/internal/app/repo/balance"
	"github.com/harunnryd/tempolalu/internal/app/repo/mutation"
	"github.com/harunnryd/tempolalu/internal/app/repo/mutationtransaction"
	"github.com/harunnryd/tempolalu/internal/app/repo/transaction"
)

type Repo interface {
	Mutation() mutation.Mutation
	Transaction() transaction.Transaction
	MutationTransaction() mutationtransaction.MutationTransaction
	Balance() balance.Balance
}

type repo struct {
	mutation            mutation.Mutation
	transaction         transaction.Transaction
	mutationtransaction mutationtransaction.MutationTransaction
	balance             balance.Balance
}

func NewRepo(cfg config.Config) Repo {
	dbase := db.New(db.WithConfig(cfg))
	pgsqlConn, err := dbase.Manager(db.PgsqlDialectParam)

	if err != nil {
		log.Fatalln("error1", err)
	}

	repo := new(repo)
	repo.mutation = mutation.New(cfg, pgsqlConn)
	repo.transaction = transaction.New(cfg, pgsqlConn)
	repo.mutationtransaction = mutationtransaction.New(cfg, pgsqlConn)
	repo.balance = balance.New(cfg, pgsqlConn)

	return repo
}

func (r *repo) Mutation() mutation.Mutation {
	return r.mutation
}

func (r *repo) Transaction() transaction.Transaction {
	return r.transaction
}

func (r *repo) MutationTransaction() mutationtransaction.MutationTransaction {
	return r.mutationtransaction
}

func (r *repo) Balance() balance.Balance {
	return r.balance
}
