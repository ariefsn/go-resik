package main

import (
	"fmt"

	"github.com/ariefsn/go-resik/app/todo/delivery/api"
	"github.com/ariefsn/go-resik/app/todo/repository/mongo"
	"github.com/ariefsn/go-resik/app/todo/service"
	"github.com/ariefsn/go-resik/common"
	"github.com/ariefsn/go-resik/helper"
	"github.com/ariefsn/go-resik/logger"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	helper.InitEnv()
	logger.InitLogger()
}

func main() {
	env := helper.Env()

	// Setup db
	dbEnv := env.Mongo
	dbAddress := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbEnv.User, dbEnv.Password, dbEnv.Host, dbEnv.Port)
	client, _ := helper.MongoClient(dbAddress)
	db := client.Database(dbEnv.Db)

	// Setup Repositories
	todoRepo := mongo.NewMongoTodoRepository(db)

	// Setup Services
	todoSvc := service.NewTodoService(todoRepo)

	// Setup Apis
	todoApi := api.NewTodoApi(todoSvc)

	app := fiber.New()

	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(helper.JsonSuccess("OK"))
	})

	v1 := app.Group("/v1")
	v1.Mount("/todos", todoApi)

	app.Use(func(c *fiber.Ctx) error {
		logger.Info("[OUTBOND]", common.M{
			"path": c.Path(),
		})
		return c.Next()
	})

	addr := fmt.Sprintf("%s:%s", env.App.Host, env.App.Port)
	app.Listen(addr)
}
