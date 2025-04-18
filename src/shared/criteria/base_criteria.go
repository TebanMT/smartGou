package criteria

import (
	"fmt"

	"github.com/TebanMT/smartGou/src/modules/reference/domain"
	sharedDomain "github.com/TebanMT/smartGou/src/shared/domain"
	"github.com/TebanMT/smartGou/src/shared/utils"
)

type BaseCriteria struct {
	Limit          *int
	Offset         *int
	OrderBy        *string
	OrderDir       *string // ASC / DESC
	Paged          *bool
	TotalOfRecords *int
}

func (q *BaseCriteria) Validate() error {
	// 1) if paged is true, validate limit & offset
	if isTrue(q.Paged) {
		if err := q.validatePaging(); err != nil {
			return err
		}
	}

	// 2) validate ordering rules
	if err := q.validateOrdering(); err != nil {
		return err
	}

	// 3) if OrderBy is provided but not OrderDir, set ASC by default
	q.applyDefaultOrderDir()

	return nil
}

func isTrue(b *bool) bool {
	return b != nil && *b
}

func (q *BaseCriteria) validatePaging() error {
	if q.Limit == nil || q.Offset == nil {
		return sharedDomain.NewValidationError(domain.ErrLimitAndOffsetRequired)
	}
	if *q.Limit < 1 {
		return sharedDomain.NewValidationError(domain.ErrLimitMustBeGreaterThanZero)
	}
	if *q.Offset < 0 {
		return sharedDomain.NewValidationError(domain.ErrOffsetMustBeGreaterThanZero)
	}
	return nil
}

func (q *BaseCriteria) validateOrdering() error {
	// if OrderDir is provided, then OrderBy is required and the direction must be valid
	if q.OrderDir != nil {
		if q.OrderBy == nil || *q.OrderBy == "" {
			return sharedDomain.NewValidationError(domain.ErrOrderByRequired)
		}
		switch *q.OrderDir {
		case "ASC", "DESC":
		default:
			return sharedDomain.NewValidationError(domain.ErrInvalidOrderDirection)
		}
	}

	// if OrderBy is provided but empty, return error
	if q.OrderBy != nil && *q.OrderBy == "" {
		return sharedDomain.NewValidationError(domain.ErrOrderByRequired)
	}

	return nil
}

func (q *BaseCriteria) applyDefaultOrderDir() {
	if q.OrderBy != nil && q.OrderDir == nil {
		dir := "ASC"
		q.OrderDir = &dir
	}
}

func (q *BaseCriteria) Debug() string {
	return fmt.Sprintf("BaseCriteria{Limit: %v, Offset: %v, OrderBy: %v, OrderDir: %v, Paged: %v}", utils.Safe(q.Limit), utils.Safe(q.Offset), utils.Safe(q.OrderBy), utils.Safe(q.OrderDir), utils.Safe(q.Paged))
}
