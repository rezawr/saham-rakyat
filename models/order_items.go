package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type OrderItem struct {
	Id             int       `json:"id" gorm:"primary_key"`
	Name           string    `json:"name"`
	Price          int       `json:"price"`
	ExpiredAt      time.Time `json:"expired_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      soft_delete.DeletedAt
	OrderHistories []OrderHistory
}
