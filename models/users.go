package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	Id             int       `json:"id" gorm:"primary_key"`
	Full_name      string    `json:"full_name"`
	First_order    string    `json:"first_order"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      soft_delete.DeletedAt
	OrderHistories []OrderHistory
}
