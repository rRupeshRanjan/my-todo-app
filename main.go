package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"my-todo-app/config"
	"my-todo-app/services"
)

func main() {
	app := fiber.New()

	defer func() { _ = app.Shutdown() }()

	configureApp(app)
	registerRoutes(app)

	err := app.Listen(config.Port)
	if err != nil {
		log.Panic("Error starting server with error: ", err)
	}
}

func configureApp(app *fiber.App) {
	app.Use(
		config.GetFiberLogger(),
		config.GetCors(),
	)
}

func registerRoutes(app *fiber.App) {
	app.Get("/task/:id", services.GetTaskByIdHandler)
	app.Get("/tasks", services.GetAllTasksHandler)
	app.Get("/tasks/search", services.SearchHandler)
	app.Post("/task", services.CreateTaskHandler)
	app.Put("/task/:id", services.UpdateTaskByIdHandler)
	app.Delete("/task/:id", services.DeleteTaskByIdHandler)
}
