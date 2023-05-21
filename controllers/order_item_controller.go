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
	"gorm.io/gorm/clause"
)

type CreateOrderItemInput struct {
	Name      string    `json:"name" binding:"required" validate:"required"`
	Price     int       `json:"price" binding:"required" validate:"required"`
	ExpiredAt time.Time `json:"expired_at" binding:"required" validate:"required"`
}

func FetchAllOrderItems(c echo.Context) error {
	var res Response

	db := c.Get("db").(*gorm.DB)
	redisClient := c.Get("redis").(*redis.Client)

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("order_items_page_%d", page)
	result, err := redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var order_items []models.OrderItem
		json.Unmarshal([]byte(result), &order_items)
		res.Status = http.StatusOK
		res.Message = "SUCCESS"
		res.Data = order_items

		return c.JSON(http.StatusOK, res)
	}

	var order_items []models.OrderItem
	limit := 10
	offset := (page - 1) * limit
	db.Offset(offset).Limit(limit).Find(&order_items)

	jsonBytes, _ := json.Marshal(order_items)
	redisClient.Set(context.Background(), cacheKey, jsonBytes, 5*time.Minute)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_items

	return c.JSON(http.StatusOK, res)
}

func FetchOrderItem(c echo.Context) error {
	var order_item models.OrderItem
	var res Response

	if err := db.DB.Where("id = ?", c.Param("id")).First(&order_item).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_item

	return c.JSON(http.StatusOK, res)
}

func CreateOrderItem(c echo.Context) error {
	var input CreateOrderItemInput
	var res Response

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

	order_item := models.OrderItem{
		Name:      input.Name,
		Price:     input.Price,
		ExpiredAt: input.ExpiredAt,
	}
	db.DB.Create(&order_item)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_item

	return c.JSON(http.StatusOK, res)
}

func UpdateOrderItem(c echo.Context) error {
	var input CreateOrderItemInput
	var res Response
	var order_item models.OrderItem

	if err := db.DB.Where("id = ?", c.Param("id")).First(&order_item).Error; err != nil {
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
	db.DB.Model(&order_item).Updates(input)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_item

	return c.JSON(http.StatusOK, res)
}

func DeleteOrderItem(c echo.Context) error {
	var order_item models.OrderItem
	var res Response

	if err := db.DB.Where("id = ?", c.Param("id")).First(&order_item).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db.DB.Select(clause.Associations).Delete(&order_item)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = order_item

	return c.JSON(http.StatusOK, res)
}
