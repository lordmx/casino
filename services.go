package main

import (
	"github.com/jinzhu/gorm"
)

// Сервис экспресс-кредитовая
type ServiceExpress struct {
	DB *gorm.DB
	Clients []Client
}
