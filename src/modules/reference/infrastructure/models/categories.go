package models

import "github.com/google/uuid"

type MetaCategory struct {
	MetaCategoryID uuid.UUID `gorm:"primaryKey;type:uuid"`
	NameEn         string    `gorm:"not null"`
	NameEs         string    `gorm:"not null"`
	Icon           string    `gorm:"not null"`
	Color          string    `gorm:"not null"`
	Description    string    `gorm:"not null"`
}

func (MetaCategory) TableName() string {
	return "meta_categories"
}

type SpendingCategory struct {
	SpendingCategoryID uuid.UUID `gorm:"primaryKey;type:uuid"`
	NameEn             string    `gorm:"not null"`
	NameEs             string    `gorm:"not null"`
	Icon               string    `gorm:"not null"`
	Color              string    `gorm:"not null"`
	Description        string    `gorm:"not null"`
	MetaCategoryID     uuid.UUID `gorm:"not null"`
	MetaCategory       MetaCategory
}

func (SpendingCategory) TableName() string {
	return "spending_categories"
}
