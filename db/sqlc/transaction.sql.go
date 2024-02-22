package db

import "context"

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (
    transaction_id,
	request_id,
	merchant,
    customer,
	amount,
	currency,
	concept
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING transaction_id, request_id, merchant, customer, amount, currency, concept, created_at
`
const getTransaction = `-- name: GetTransaction :one
SELECT * FROM transactions WHERE transaction_id = $1 LIMIT 1

`

type CreateTransactionHistory struct {
	TransactionId string  `json:"transaction_id"`
	RequestId     string  `json:"request_id"`
	Merchant      string  `json:"merchant"`
	Customer      string  `json:"customer"`
	Amount        float64 `json:"amount"`
	Concept       string  `json:"concept"`
	Currency      string  `json:"currency"`
}

func (q *Queries) CreateTransactionHistory(ctx context.Context, arg CreateTransactionHistory) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.TransactionId,
		arg.RequestId,
		arg.Merchant,
		arg.Customer,
		arg.Amount,
		arg.Currency,
		arg.Concept,
	)

	var t Transaction
	err := row.Scan(
		&t.TransactionId,
		&t.RequestId,
		&t.Merchant,
		&t.Customer,
		&t.Amount,
		&t.Currency,
		&t.Concept,
		&t.CreatedAt,
	)
	return t, err
}

func (q *Queries) GetTransactionDetails(ctx context.Context, transactionId string) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, transactionId)
	var t Transaction
	err := row.Scan(
		&t.TransactionId,
		&t.RequestId,
		&t.Merchant,
		&t.Customer,
		&t.Amount,
		&t.Currency,
		&t.Concept,
		&t.CreatedAt,
	)
	return t, err
}
