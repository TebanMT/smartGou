package domain

import (
	"context"

	"gorm.io/gorm"
)

// GormTransaction is a concrete implementation of the Transaction interface
type GormTransaction struct {
	Tx *gorm.DB
}

func (g *GormTransaction) Execute(fn func(tx Transaction) error) error {
	return fn(g)
}

// gormUnitOfWork is a concrete implementation of the UnitOfWork interface
type gormUnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &gormUnitOfWork{db: db}
}

func (u *gormUnitOfWork) Begin(ctx context.Context) (Transaction, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormTransaction{Tx: tx}, nil
}

func (u *gormUnitOfWork) Commit(tx Transaction) error {
	return tx.(*GormTransaction).Tx.Commit().Error
}

func (u *gormUnitOfWork) Rollback(tx Transaction) error {
	return tx.(*GormTransaction).Tx.Rollback().Error
}

func (u *gormUnitOfWork) Query(ctx context.Context) (Transaction, error) {
	tx := u.db.WithContext(ctx)
	return &GormTransaction{Tx: tx}, nil
}

func (u *gormUnitOfWork) Command(ctx context.Context, fn func(tx Transaction) error) error {
	tx, err := u.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			u.Rollback(tx)
			panic(r) // volver a lanzar el panic
		}
	}()

	if err := fn(tx); err != nil {
		u.Rollback(tx)
		return err
	}

	return u.Commit(tx)
}
