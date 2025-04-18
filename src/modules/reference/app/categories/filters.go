package app

import (
	"fmt"

	"github.com/TebanMT/smartGou/src/modules/reference/domain"
	"github.com/TebanMT/smartGou/src/shared/criteria"
	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/google/uuid"
)

type CategoryQuery struct {
	MetaCategoryID *uuid.UUID
	NameLike       *string
	criteria.BaseCriteria
}

func (q *CategoryQuery) Validate() error {
	err := q.BaseCriteria.Validate()
	fmt.Println("UU", q.MetaCategoryID)
	if err != nil {
		return err
	}
	if q.MetaCategoryID != nil {
		if _, err := uuid.Parse(q.MetaCategoryID.String()); err != nil {
			return sharedDomain.NewValidationError(domain.ErrInvalidMetaCategoryID)
		}
	}
	return nil
}
