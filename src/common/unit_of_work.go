package common

import (
	"context"

	"github.com/TebanMT/smartGou/src/common/domain"
	"gorm.io/gorm"
)

type GormTransaction struct {
	Tx *gorm.DB
}

func (g *GormTransaction) Execute(fn func(tx domain.Transaction) error) error {
	return fn(g)
}

type UnitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) domain.UnitOfWork {
	return &UnitOfWork{db: db}
}

func (u *UnitOfWork) Begin(ctx context.Context) (domain.Transaction, error) {
	tx := u.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &GormTransaction{Tx: tx}, nil
}

func (u *UnitOfWork) Commit(tx domain.Transaction) error {
	return tx.(*GormTransaction).Tx.Commit().Error
}

func (u *UnitOfWork) Rollback(tx domain.Transaction) error {
	return tx.(*GormTransaction).Tx.Rollback().Error
}
