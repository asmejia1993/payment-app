package db

import "context"

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUser(ctx context.Context, email string) (User, error)
	CreateTransactionHistory(ctx context.Context, arg CreateTransactionHistory) (Transaction, error)
	GetTransactionDetails(ctx context.Context, transactionId string) (Transaction, error)
}

var _ Querier = (*Queries)(nil)
