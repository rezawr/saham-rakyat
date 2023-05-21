package db

import (
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rezawr/saham-rakyat/config"
	"github.com/rezawr/saham-rakyat/models"
)

var err error
var DB *gorm.DB
var RedisClient *redis.Client

func Init() {
	conf := config.GetConf()

	dsn := conf.DB_USERNAME + ":" + conf.DB_PASSWORD + "@tcp(" + conf.DB_HOST + ":" + conf.DB_PORT + ")/" + conf.DB_NAME + "?parseTime=true"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Connection error")
	}

	// Migrate the schema
	database.AutoMigrate(&models.User{}, &models.OrderItem{}, &models.OrderHistory{})

	DB = database
}

func InitRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	RedisClient = client
	return client
}
