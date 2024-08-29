package models

type IngredientsRecipes struct {
	RecipeID     uint   `gorm:"primaryKey"`
	IngredientID uint   `gorm:"primaryKey"`
	Quantity     string `gorm:"not null"`
	Recipe       uint   `gorm:"foreignKey:RecipeID"`
	Ingredient   uint   `gorm:"foreignKey:IngredientID"`
}
