package app

import (
	"context"

	"github.com/TebanMT/smartGou/src/modules/reference/domain"
	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
)

type GetCategoriesUseCase struct {
	categoryRepository domain.CategoriesRepository
	unitOfWork         sharedDomain.UnitOfWork
}

func NewGetCategoriesUseCase(categoryRepository domain.CategoriesRepository, unitOfWork sharedDomain.UnitOfWork) *GetCategoriesUseCase {
	return &GetCategoriesUseCase{categoryRepository: categoryRepository, unitOfWork: unitOfWork}
}

func (u *GetCategoriesUseCase) GetCategories(ctx context.Context, criteria *CategoryQuery) ([]sharedDomain.Category, error) {
	if err := criteria.Validate(); err != nil {
		return nil, err
	}
	tx, err := u.unitOfWork.Query(ctx)
	if err != nil {
		return nil, err
	}
	if criteria.Paged != nil && *criteria.Paged {
		total, err := u.categoryRepository.CountCategories(tx, criteria)
		if err != nil {
			return nil, err
		}
		criteria.TotalOfRecords = &total
	}
	categories, err := u.categoryRepository.FindCategories(tx, criteria)
	if err != nil {
		return nil, err
	}
	return categories, nil
}
