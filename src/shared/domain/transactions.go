package domain

import (
	"context"
)

type Transaction interface {
	Execute(fn func(tx Transaction) error) error
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
	Commit(tx Transaction) error
	Rollback(tx Transaction) error
	Command(ctx context.Context, fn func(tx Transaction) error) error
	// Query is used to get a transaction without starting a new one. It is used to execute queries that don't need to be part of a transaction.
	// Is equivalent to 'Begin' but is used to be more explicit.
	Query(ctx context.Context) (Transaction, error)
}
