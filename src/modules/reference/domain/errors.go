package domain

import "errors"

var (
	ErrLimitAndOffsetRequired      = errors.New("limit and offset are required when paged is true")
	ErrLimitGreaterThanOffset      = errors.New("limit must not be greater than offset")
	ErrInvalidPaginationRange      = errors.New("invalid pagination range")
	ErrUnexpectedPaginationParams  = errors.New("unexpected pagination params. limit and offset must be provided only when paged is true")
	ErrLimitMustBeGreaterThanZero  = errors.New("limit must be greater than 0")
	ErrOffsetMustBeGreaterThanZero = errors.New("offset must be greater than 0")
	ErrOrderByRequired             = errors.New("order by is required when direction is provided")
	ErrInvalidOrderDirection       = errors.New("invalid order direction DESC or ASC")
	ErrInvalidMetaCategoryID       = errors.New("invalid meta category id")
)
