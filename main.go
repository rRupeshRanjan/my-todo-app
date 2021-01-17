package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"my-todo-app/config"
	"my-todo-app/services"
)

func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: config.FiberLogFormat,
		Output: config.LogFile,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	registerRoutes(app)

	err := app.Listen(config.Port)
	if err != nil {
		log.Panic("Error starting server with error: ", err)
	}
}

func registerRoutes(app *fiber.App) {
	app.Get("/task/:id", services.GetTaskByIdHandler)
	app.Get("/tasks", services.GetAllTasksHandler)
	app.Get("/tasks/search", services.SearchHandler)
	app.Post("/task", services.CreateTaskHandler)
	app.Put("/task/:id", services.UpdateTaskByIdHandler)
	app.Put("/tasks/bulk-action", services.UpdateBulkTaskHandler)
	app.Delete("/task/:id", services.DeleteTaskByIdHandler)
}
