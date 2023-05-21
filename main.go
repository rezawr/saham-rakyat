package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rezawr/saham-rakyat/db"
	"github.com/rezawr/saham-rakyat/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("ENV Not working")
	}

	db.Init()
	db.InitRedisClient()

	e := routes.Init()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db.DB)
			c.Set("redis", db.RedisClient)
			return next(c)
		}
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
