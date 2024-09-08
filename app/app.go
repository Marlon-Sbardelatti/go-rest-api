package app

import (
	"gorm.io/gorm"
)

// Objeto de acesso aos dados (DAO), que intermedia a interação com o banco
type App struct {
	DB *gorm.DB
}
