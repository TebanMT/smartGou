package criteria

import (
	"fmt"

	"gorm.io/gorm"
)

type GormCriteriaConverter struct{}

func (c *GormCriteriaConverter) ConvertBaseCriteriaToQuery(query *gorm.DB, criteria BaseCriteria) (*gorm.DB, error) {
	if criteria.Paged != nil && *criteria.Paged {
		query = query.Limit(*criteria.Limit)
		query = query.Offset(*criteria.Offset)
	}
	if criteria.OrderBy != nil && *criteria.OrderBy != "" {
		orderInstruction := fmt.Sprintf("%s %s", *criteria.OrderBy, *criteria.OrderDir)
		query = query.Order(orderInstruction)
	}

	return query, nil
}
