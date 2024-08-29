package models

type Recipe struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null"`
	Name         string `gorm:"unique;not null"`
	Instructions string `gorm:"not null"`
	User         User   `gorm:"foreignKey:UserID"`
}
