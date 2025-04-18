package domain

import (
	"github.com/TebanMT/smartGou/src/shared/domain"
)

type CategoriesRepository interface {
	FindMetaCategories(tx domain.Transaction, criteria domain.Criteria) ([]domain.MetaCategory, error)
	FindCategories(tx domain.Transaction, criteria domain.Criteria) ([]domain.Category, error)
	CountMetaCategories(tx domain.Transaction, criteria domain.Criteria) (int, error)
	CountCategories(tx domain.Transaction, criteria domain.Criteria) (int, error)
}
