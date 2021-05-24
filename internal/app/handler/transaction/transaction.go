package transaction

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/harunnryd/tempolalu/internal/app/usecase"
)

type Transaction interface {
	CreateTransaction(w http.ResponseWriter, r *http.Request) (resp interface{}, err error)
	DeleteTransactionByID(w http.ResponseWriter, r *http.Request) (resp interface{}, err error)
}

type transaction struct {
	usecase usecase.Usecase
}

func New(usecase usecase.Usecase) Transaction {
	return &transaction{usecase: usecase}
}

type CreateTransaction struct {
	OriginID      int     `json:"origin_id"`
	DestinationID int     `json:"destination_id"`
	ProductID     int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Amount        float64 `json:"amount"`
}

func (t *transaction) CreateTransaction(w http.ResponseWriter, r *http.Request) (resp interface{}, err error) {
	var createTransaction CreateTransaction
	if err = json.NewDecoder(r.Body).Decode(&createTransaction); err != nil {
		return
	}

	resp, err = t.usecase.Transaction().CreateTransaction(
		r.Context(),
		createTransaction.OriginID,
		createTransaction.DestinationID,
		createTransaction.ProductID,
		createTransaction.Quantity,
		createTransaction.Amount,
	)

	return
}

func (t *transaction) DeleteTransactionByID(w http.ResponseWriter, r *http.Request) (resp interface{}, err error) {
	transactionID, err := strconv.Atoi(chi.URLParam(r, "transaction_id"))
	if err != nil {
		return
	}

	err = t.usecase.Transaction().DeleteTransactionByID(r.Context(), transactionID)

	return
}
