package repositories

import (
	app "github.com/TebanMT/smartGou/src/modules/reference/app/categories"
	"github.com/TebanMT/smartGou/src/modules/reference/domain"
	"github.com/TebanMT/smartGou/src/modules/reference/infrastructure/models"
	criteria "github.com/TebanMT/smartGou/src/shared/criteria"
	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"gorm.io/gorm"
)

type GormCategoriesRepository struct {
	db        *gorm.DB
	converter *criteria.GormCriteriaConverter
}

func NewGormCategoriesRepository(db *gorm.DB) domain.CategoriesRepository {
	return &GormCategoriesRepository{db: db, converter: &criteria.GormCriteriaConverter{}}
}

func (r *GormCategoriesRepository) FindMetaCategories(tx sharedDomain.Transaction, c sharedDomain.Criteria) ([]sharedDomain.MetaCategory, error) {
	var metaCategories []models.MetaCategory
	err := tx.Execute(func(tx sharedDomain.Transaction) error {
		gormTx := tx.(*sharedDomain.GormTransaction)
		filterCriteria := c.(*app.CategoryQuery)
		query, err := r.converter.ConvertBaseCriteriaToQuery(gormTx.Tx, filterCriteria.BaseCriteria)
		if err != nil {
			return err
		}
		if filterCriteria.NameLike != nil {
			query = query.Where("LOWER(name_en) LIKE LOWER(?)", "%"+*filterCriteria.NameLike+"%")
		}
		return query.Find(&metaCategories).Error
	})
	if err != nil {
		return nil, err
	}
	response := []sharedDomain.MetaCategory{}
	for _, metaCategory := range metaCategories {
		response = append(response, sharedDomain.MetaCategory{
			ID:          metaCategory.MetaCategoryID,
			NameEn:      metaCategory.NameEn,
			NameEs:      metaCategory.NameEs,
			Icon:        metaCategory.Icon,
			Color:       metaCategory.Color,
			Description: metaCategory.Description,
		})
	}
	return response, nil
}

func (r *GormCategoriesRepository) CountMetaCategories(tx sharedDomain.Transaction, c sharedDomain.Criteria) (int, error) {
	var count int64
	err := tx.Execute(func(tx sharedDomain.Transaction) error {
		gormTx := tx.(*sharedDomain.GormTransaction)
		filterCriteria := c.(*app.CategoryQuery)
		query := gormTx.Tx.Model(&models.MetaCategory{})
		if filterCriteria.NameLike != nil {
			query = query.Where("LOWER(name_en) LIKE LOWER(?)", "%"+*filterCriteria.NameLike+"%")
		}
		return query.Count(&count).Error
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GormCategoriesRepository) CountCategories(tx sharedDomain.Transaction, c sharedDomain.Criteria) (int, error) {
	var count int64
	err := tx.Execute(func(tx sharedDomain.Transaction) error {
		gormTx := tx.(*sharedDomain.GormTransaction)
		filterCriteria := c.(*app.CategoryQuery)
		query := gormTx.Tx.Model(&models.SpendingCategory{})
		if filterCriteria.MetaCategoryID != nil {
			query = query.Where("meta_category_id = ?", filterCriteria.MetaCategoryID)
		}
		if filterCriteria.NameLike != nil {
			query = query.Where("LOWER(name_en) LIKE LOWER(?)", "%"+*filterCriteria.NameLike+"%")
		}
		return query.Count(&count).Error
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *GormCategoriesRepository) FindCategories(tx sharedDomain.Transaction, c sharedDomain.Criteria) ([]sharedDomain.Category, error) {
	var categories []models.SpendingCategory
	err := tx.Execute(func(tx sharedDomain.Transaction) error {
		gormTx := tx.(*sharedDomain.GormTransaction)
		filterCriteria := c.(*app.CategoryQuery)
		query, err := r.converter.ConvertBaseCriteriaToQuery(gormTx.Tx, filterCriteria.BaseCriteria)
		if err != nil {
			return err
		}
		if filterCriteria.MetaCategoryID != nil {
			query = query.Where("meta_category_id = ?", filterCriteria.MetaCategoryID)
		}
		if filterCriteria.NameLike != nil {
			query = query.Where("LOWER(name_en) LIKE LOWER(?)", "%"+*filterCriteria.NameLike+"%")
		}
		return query.Find(&categories).Error
	})
	if err != nil {
		return nil, err
	}
	response := []sharedDomain.Category{}
	for _, category := range categories {
		response = append(response, sharedDomain.Category{
			ID:             category.SpendingCategoryID,
			NameEn:         category.NameEn,
			NameEs:         category.NameEs,
			Icon:           category.Icon,
			Color:          category.Color,
			Description:    category.Description,
			MetaCategoryID: category.MetaCategoryID.String(),
		})
	}
	return response, nil
}
