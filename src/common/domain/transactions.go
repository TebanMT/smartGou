package domain

import "context"

type Transaction interface {
	Execute(fn func(tx Transaction) error) error
}

type UnitOfWork interface {
	Begin(ctx context.Context) (Transaction, error)
	Commit(tx Transaction) error
	Rollback(tx Transaction) error
}
