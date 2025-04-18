package domain

import "github.com/google/uuid"

type MetaCategory struct {
	ID          uuid.UUID
	NameEn      string
	NameEs      string
	Icon        string
	Color       string
	Description string
}

type Category struct {
	ID             uuid.UUID
	NameEn         string
	NameEs         string
	Icon           string
	Color          string
	Description    string
	MetaCategoryID string
	MetaCategory   MetaCategory
}
