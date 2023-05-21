package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/rezawr/saham-rakyat/db"
	"github.com/rezawr/saham-rakyat/models"
	"gorm.io/gorm"
)

type CreateOrderHistoryInput struct {
	UserId      int    `json:"user_id" binding:"required" validate:"required"`
	OrderItemId int    `json:"order_item_id" binding:"required" validate:"required"`
	Description string `json:"description" binding:"required"`
}

func FetchAllOrderHistories(c echo.Context) error {
	var res Response

	db := c.Get("db").(*gorm.DB)
	redisClient := c.Get("redis").(*redis.Client)

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("order_histories_page_%d", page)
	result, err := redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var order_histories []models.OrderHistory
		json.Unmarshal([]byte(result), &order_histories)
		res.Status = http.StatusOK
		res.Message = "SUCCESS"
		res.Data = order_histories

		return c.JSON(http.StatusOK, res)
	}

	var order_histories []models.OrderHistory
	limit := 10
	offset := (page - 1) * limit
	db.Preload("User").Preload("OrderItem").Offset(offset).Limit(limit).Find(&order_histories)

	jsonBytes, _ := json.Marshal(order_histories)
	redisClient.Set(context.Background(), cacheKey, jsonBytes, 5*time.Minute)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_histories

	return c.JSON(http.StatusOK, res)
}

func FetchOrderHistory(c echo.Context) error {
	var order_history models.OrderHistory
	var res Response

	if err := db.DB.Where("id = ?", c.Param("id")).Preload("User").Preload("OrderItem").First(&order_history).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_history

	return c.JSON(http.StatusOK, res)
}

func CreateOrderHistory(c echo.Context) error {
	var input CreateOrderHistoryInput
	var res Response
	var user models.User
	var order_item models.OrderItem

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	v := validator.New()

	err := v.Struct(input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	if err := db.DB.Where("id = ?", input.UserId).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "User not found"})
	}

	if err := db.DB.Where("id = ?", input.OrderItemId).First(&order_item).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Order Item not found"})
	}

	order_history := models.OrderHistory{
		UserId:      input.UserId,
		User:        user,
		OrderItemId: input.OrderItemId,
		OrderItem:   order_item,
		Description: input.Description,
	}
	db.DB.Create(&order_history)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_history

	return c.JSON(http.StatusOK, res)
}

func UpdateOrderHistory(c echo.Context) error {
	var input CreateOrderHistoryInput
	var res Response
	var order_history models.OrderHistory

	if err := db.DB.Where("id = ?", c.Param("id")).Preload("User").Preload("OrderItem").First(&order_history).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	v := validator.New()

	err := v.Struct(input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	db.DB.Model(&order_history).Updates(input)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_history

	return c.JSON(http.StatusOK, res)
}

func DeleteOrderHistory(c echo.Context) error {
	var order_history models.OrderHistory
	var res Response

	if err := db.DB.Where("id = ?", c.Param("id")).First(&order_history).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db.DB.Delete(&order_history)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_history

	return c.JSON(http.StatusOK, res)
}
