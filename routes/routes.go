package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rezawr/saham-rakyat/controllers"
)

func Init() *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/users", controllers.FetchAllUsers)
	e.GET("/users/:id", controllers.FetchUser)
	e.POST("/users", controllers.CreateUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	e.GET("/order_items", controllers.FetchAllOrderItems)
	e.GET("/order_items/:id", controllers.FetchOrderItem)
	e.POST("/order_items", controllers.CreateOrderItem)
	e.PUT("/order_items/:id", controllers.UpdateOrderItem)
	e.DELETE("/order_items/:id", controllers.DeleteOrderItem)

	e.GET("/order_histories", controllers.FetchAllOrderHistories)
	e.GET("/order_histories/:id", controllers.FetchOrderHistory)
	e.POST("/order_histories", controllers.CreateOrderHistory)
	e.PUT("/order_histories/:id", controllers.UpdateOrderHistory)
	e.DELETE("/order_histories/:id", controllers.DeleteOrderHistory)

	return e
}
