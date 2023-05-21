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

type CreateUserInput struct {
	Full_name   string `json:"full_name" binding:"required" validate:"required"`
	First_order string `json:"first_order" binding:"required" validate:"required"`
}

func FetchAllUsers(c echo.Context) error {
	var res Response

	db := c.Get("db").(*gorm.DB)
	redisClient := c.Get("redis").(*redis.Client)

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	cacheKey := fmt.Sprintf("users_page_%d", page)
	result, err := redisClient.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var users []models.User
		json.Unmarshal([]byte(result), &users)
		res.Status = http.StatusOK
		res.Message = "SUCCESS"
		res.Data = users

		return c.JSON(http.StatusOK, res)
	}

	var users []models.User
	limit := 10
	offset := (page - 1) * limit
	db.Offset(offset).Limit(limit).Find(&users)

	jsonBytes, _ := json.Marshal(users)
	redisClient.Set(context.Background(), cacheKey, jsonBytes, 5*time.Minute)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = users

	return c.JSON(http.StatusOK, res)
}

func FetchUser(c echo.Context) error {
	var user models.User
	var res Response

	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = user

	return c.JSON(http.StatusOK, res)
}

func CreateUser(c echo.Context) error {
	var input CreateUserInput
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

	user := models.User{
		Full_name:   input.Full_name,
		First_order: input.First_order,
	}
	db.DB.Create(&user)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = user

	return c.JSON(http.StatusOK, res)
}

func UpdateUser(c echo.Context) error {
	var input CreateUserInput
	var res Response
	var user models.User

	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
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
	db.DB.Model(&user).Updates(input)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = user

	return c.JSON(http.StatusOK, res)
}

func DeleteUser(c echo.Context) error {
	var user models.User
	var res Response

	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db.DB.Select(clause.Associations).Delete(&user)

	res.Status = http.StatusOK
	res.Message = "SUCCESS"
	res.Data = user

	return c.JSON(http.StatusOK, res)
}
