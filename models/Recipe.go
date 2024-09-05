package models

type Recipe struct {
	ID                 uint                 `gorm:"primaryKey" json:"id"`
	UserID             uint                 `gorm:"not null" json:"user_id"`
	Name               string               `gorm:"unique;not null" json:"name"`
	Instructions       string               `gorm:"not null" json:"instructions"`
	IngredientsRecipes []IngredientsRecipes `gorm:"foreignKey:RecipeID;constraint:OnDelete:CASCADE" json:"ingredients"`
}
