package models

type IngredientsRecipes struct {
	RecipeID     uint       `gorm:"primaryKey" json:"-"`
	IngredientID uint       `gorm:"primaryKey" json:"ingredient_id"`
	Quantity     string     `gorm:"not null" json:"quantity"`
	Ingredient   Ingredient `gorm:"foreignKey:IngredientID;constraint:OnDelete:CASCADE" json:"ingredient"`
}
