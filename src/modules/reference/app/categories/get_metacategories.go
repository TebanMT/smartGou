package app

import (
	"context"

	"github.com/TebanMT/smartGou/src/modules/reference/domain"
	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
)

type GetMetaCategoriesUseCase struct {
	categoriesRepository domain.CategoriesRepository
	unitOfWork           sharedDomain.UnitOfWork
}

func NewGetMetaCategoriesUseCase(categoriesRepository domain.CategoriesRepository, unitOfWork sharedDomain.UnitOfWork) *GetMetaCategoriesUseCase {
	return &GetMetaCategoriesUseCase{categoriesRepository: categoriesRepository, unitOfWork: unitOfWork}
}

func (u *GetMetaCategoriesUseCase) GetMetaCategories(ctx context.Context, criteria *CategoryQuery) ([]sharedDomain.MetaCategory, error) {
	if err := criteria.Validate(); err != nil {
		return nil, err
	}
	tx, err := u.unitOfWork.Query(ctx)
	if err != nil {
		return nil, err
	}

	if criteria.Paged != nil && *criteria.Paged {
		total, err := u.categoriesRepository.CountMetaCategories(tx, criteria)
		if err != nil {
			return nil, err
		}
		criteria.TotalOfRecords = &total
	}

	metaCategories, err := u.categoriesRepository.FindMetaCategories(tx, criteria)
	if err != nil {
		return nil, err
	}

	return metaCategories, nil
}
