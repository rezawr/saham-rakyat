package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type OrderHistory struct {
	Id          int `json:"id" gorm:"primary_key"`
	UserId      int
	User        User `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignkey:UserId"`
	OrderItemId int
	OrderItem   OrderItem `json:"order_item" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignkey:OrderItemId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   soft_delete.DeletedAt
}
