package model

type MutationTransaction struct {
	Model
	MutationID    int `json:"mutation_id"`
	TransactionID int `json:"transaction_id"`
}

func (MutationTransaction) TableName() string {
	return "mutation_transaction"
}
