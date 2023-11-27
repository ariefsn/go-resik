package main

import (
	"fmt"
	"net/http"

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
	client, _ := helper.MongoClient()
	db := client.Database("resik-arch")

	// Setup Repositories
	todoRepo := mongo.NewMongoTodoRepository(db)

	// Setup Services
	todoSvc := service.NewTodoService(todoRepo)

	// Setup Apis
	todoApi := api.NewTodoApi(todoSvc)

	app := fiber.New()

	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		Done: func(c *fiber.Ctx, logString []byte) {
			// fmt.Printf("[OUTBOUND] %s", string(logString))
			// logger.Info("[OUTBOUND]", common.M{
			// 	"path":   c.Path(),
			// 	"method": c.Method(),
			// 	"host":   string(c.Request().Host()),
			// 	"status": c.Response().StatusCode(),
			// 	"time":   time.Now(),
			// })
		},
	}))

	// app.Use(func(c *fiber.Ctx) error {
	// 	logger.Info("[INBOUND]", common.M{
	// 		"path": c.Path(),
	// 	})
	// 	return c.Next()
	// })

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

	methods := []string{
		http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPost,
	}
	for _, r := range app.GetRoutes() {
	MTD:
		for _, m := range methods {
			if r.Method == m {
				logger.Info("", common.M{
					"method": r.Method,
					"path":   r.Path,
				})
				break MTD
			}
		}
	}

	addr := fmt.Sprintf("%s:%s", env.App.Host, env.App.Port)
	app.Listen(addr)
}
