package services

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"my-todo-app/domain"
	"my-todo-app/testUtils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkInit(b *testing.B) {
	taskRepository = taskRepositoryMock{}
	testApp = fiber.New()
}

func BenchmarkGetTaskByIdHandler(b *testing.B) {
	testApp.Get("/task/:id", func(c *fiber.Ctx) error {
		return GetTaskByIdHandler(c)
	})
	scenarios := testUtils.GetServiceTestScenarios(testUtils.GetTaskByIdKey)

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			taskRepositoryGetByIdMock = func(id string) ([]domain.Task, error) {
				return scenario.ExpectedTasks, scenario.ScenarioErr
			}

			request := httptest.NewRequest("GET", "http://localhost.com/task/8", nil)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkGetAllTasksHandler(b *testing.B) {
	scenarios := testUtils.GetServiceTestScenarios(testUtils.GetAllTasksKey)

	testApp.Get("/expectedTasks", func(c *fiber.Ctx) error {
		return GetAllTasksHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			taskRepositoryGetAllTasksMock = func(page int64, perPage int64) ([]domain.Task, error) {
				return scenario.ExpectedTasks, scenario.ScenarioErr
			}
			request := httptest.NewRequest("GET", "http://localhost.com/expectedTasks", nil)

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkCreateTaskHandler(b *testing.B) {
	scenarios := testUtils.GetServiceTestScenarios(testUtils.CreateTaskKey)
	testApp.Post("/task", func(c *fiber.Ctx) error {
		return CreateTaskHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {

			taskRepositoryCreateTaskMock = func(task domain.Task) (int64, error) {
				return 1, scenario.ScenarioErr
			}

			request := httptest.NewRequest("POST", "http://localhost.com/task", bytes.NewBuffer(scenario.Data))
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkUpdateTaskByIdHandler(b *testing.B) {
	scenarios := testUtils.GetServiceTestScenarios(testUtils.UpdateTaskKey)

	testApp.Put("/task/:id", func(c *fiber.Ctx) error {
		return UpdateTaskByIdHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			taskRepositoryUpdateTaskMock = func(task domain.Task, id string) error {
				return scenario.ScenarioErr
			}

			request := httptest.NewRequest("PUT", "http://localhost.com/task/1", bytes.NewBuffer(scenario.Data))
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkDeleteTaskByIdHandler(b *testing.B) {
	testApp.Delete("/task/:id", func(c *fiber.Ctx) error {
		return DeleteTaskByIdHandler(c)
	})
	scenarios := testUtils.GetServiceTestScenarios(testUtils.DeleteTaskKey)

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			request := httptest.NewRequest("DELETE", "http://localhost.com/task/8", nil)
			taskRepositoryDeleteTaskMock = func(id string) (int64, error) {
				return scenario.RowsAffected, scenario.ScenarioErr
			}

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkSearchHandler(b *testing.B) {
	scenarios := testUtils.GetServiceTestScenarios(testUtils.SearchTaskKey)

	testApp.Get("/expectedTasks/search", func(c *fiber.Ctx) error {
		return SearchHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			request := httptest.NewRequest("GET", "http://localhost.com/expectedTasks/search", nil)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkUpdateBulkTaskHandler(b *testing.B) {
	scenarios := testUtils.GetServiceTestScenarios(testUtils.BulkUpdateTaskKey)

	testApp.Put("/expectedTasks/bulk-action", func(c *fiber.Ctx) error {
		return UpdateBulkTaskHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			request := httptest.NewRequest("PUT", "http://localhost.com/expectedTasks/bulk-action", nil)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func compareStatusCodes(b *testing.B, response *http.Response, s domain.Scenario) {
	if response.StatusCode != s.StatusCode {
		b.Fatalf("Expected status code: %d, instead got: %d", s.StatusCode, response.StatusCode)
	}
}
