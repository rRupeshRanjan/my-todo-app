package services

import (
	"bytes"
	"errors"
	"github.com/gofiber/fiber/v2"
	"my-todo-app/domain"
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
	scenarios := []scenario{
		{
			name: "should successfully get task by id",
			tasks: []domain.Task{{
				Id:          8,
				AddedOn:     123456789,
				DueBy:       123456789,
				Title:       "sample title",
				Description: "sample description",
				Status:      "sample status",
			}},
			statusCode: http.StatusOK,
		},
		{
			name:       "should give 404 for get task by id",
			tasks:      []domain.Task{},
			statusCode: http.StatusNotFound,
		},
		{
			name:       "should give 500 for get task by id for database errors",
			tasks:      []domain.Task{},
			err:        errors.New("error while fetching data"),
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			taskRepositoryGetByIdMock = func(id string) ([]domain.Task, error) {
				return scenario.tasks, scenario.err
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
	scenarios := []scenario{
		{
			name: "should successfully get all tasks",
			tasks: []domain.Task{
				{
					Id:          8,
					AddedOn:     123456789,
					DueBy:       123456789,
					Title:       "sample title",
					Description: "sample description",
					Status:      "sample status",
				},
				{
					Id:          9,
					AddedOn:     12345678,
					DueBy:       12345678,
					Title:       "sample title 2",
					Description: "sample description 2",
					Status:      "sample status 2",
				}},
			statusCode: http.StatusOK,
		},
		{
			name:       "should give zero tasks with status code 200",
			tasks:      []domain.Task{},
			statusCode: http.StatusOK,
		},
		{
			name:       "should give 500 for get task by id for database errors",
			tasks:      []domain.Task{},
			err:        errors.New("error while fetching data"),
			statusCode: http.StatusInternalServerError,
		},
	}

	testApp.Get("/tasks", func(c *fiber.Ctx) error {
		return GetAllTasksHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			taskRepositoryGetAllTasksMock = func() ([]domain.Task, error) {
				return scenario.tasks, scenario.err
			}
			request := httptest.NewRequest("GET", "http://localhost.com/tasks", nil)

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
	scenarios := []scenario{
		{
			name: "should successfully create task",
			task: domain.Task{
				Id:          1,
				AddedOn:     123,
				DueBy:       123,
				Title:       "sample",
				Description: "sample",
				Status:      "sample",
			},
			data:       []byte(`{"addedOn": 123, "dueBy": 123, "title": "sample", "description": "sample", "status": "sample"}`),
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name: "should create task overriding the id from request body",
			task: domain.Task{
				Id:          1,
				AddedOn:     123,
				DueBy:       123,
				Title:       "sample",
				Description: "sample",
				Status:      "sample",
			},
			data:       []byte(`{"id": 8, "addedOn": 123, "dueBy": 123, "title": "sample", "description": "sample", "status": "sample"}`),
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name:       "should throw 400 in create task for malformed body",
			data:       []byte(`{addedOn": 123456789, "dueBy": 123456789, "title": "sample`),
			statusCode: http.StatusBadRequest,
			err:        nil,
		},
		{
			name:       "should throw 500 in create task for database errors",
			data:       []byte(`{"addedOn": 123, "dueBy": 123, "title": "sample", "description": "sample", "status": "sample"}`),
			statusCode: http.StatusInternalServerError,
			err:        errors.New("error creating task in database"),
		},
	}
	testApp.Post("/task", func(c *fiber.Ctx) error {
		return CreateTaskHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {

			taskRepositoryCreateTaskMock = func(task domain.Task) (int64, error) {
				return 1, scenario.err
			}

			request := httptest.NewRequest("POST", "http://localhost.com/task", bytes.NewBuffer(scenario.data))
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
	scenarios := []scenario{
		{
			name: "should successfully update a task",
			task: domain.Task{
				Id:          1,
				AddedOn:     123,
				DueBy:       123,
				Title:       "sample",
				Description: "sample",
				Status:      "sample",
			},
			data:       []byte(`{"id": 1, "addedOn": 123, "dueBy": 123, "title": "sample", "description": "sample", "status": "sample"}`),
			statusCode: http.StatusOK,
			err:        nil,
		},
		{
			name:       "should throw 400 in update task if IDs are different in URL and request body",
			data:       []byte(`{"id": 8, "addedOn": 123, "dueBy": 123, "title": "sample", "description": "sample", "status": "sample"}`),
			statusCode: http.StatusBadRequest,
			err:        nil,
		},
		{
			name:       "should throw 400 in update task for malformed body",
			data:       []byte(`{"id": 1, "addedOn": 123, "dueBy": 123, "title": "sample}`),
			statusCode: http.StatusBadRequest,
			err:        nil,
		},
		{
			name:       "should throw 500 in update task database errors",
			data:       []byte(`{"id": 1, "addedOn": 123, "dueBy": 123, "title": "sample"}`),
			statusCode: http.StatusInternalServerError,
			err:        errors.New("error occurred updating task in database"),
		},
	}

	testApp.Put("/task/:id", func(c *fiber.Ctx) error {
		return UpdateTaskByIdHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			taskRepositoryUpdateTaskMock = func(task domain.Task, id string) error {
				return scenario.err
			}

			request := httptest.NewRequest("PUT", "http://localhost.com/task/1", bytes.NewBuffer(scenario.data))
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
	scenarios := []scenario{
		{
			name:         "should successfully delete task",
			statusCode:   http.StatusNoContent,
			err:          nil,
			rowsAffected: 1,
		},
		{
			name:         "should throw 404 delete task if not present in database",
			statusCode:   http.StatusNotFound,
			err:          nil,
			rowsAffected: 0,
		},
		{
			name:       "should throw 500 in delete task for database errors",
			statusCode: http.StatusInternalServerError,
			err:        errors.New("error deleting record from database"),
		},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			request := httptest.NewRequest("DELETE", "http://localhost.com/task/8", nil)
			taskRepositoryDeleteTaskMock = func(id string) (int64, error) {
				return scenario.rowsAffected, scenario.err
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
	scenarios := []scenario{
		{
			name:       "should throw 501",
			statusCode: http.StatusNotImplemented,
		},
	}

	testApp.Get("/tasks/search", func(c *fiber.Ctx) error {
		return SearchHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			request := httptest.NewRequest("GET", "http://localhost.com/tasks/search", nil)
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
	scenarios := []scenario{
		{
			name:       "should throw 501",
			statusCode: http.StatusNotImplemented,
		},
	}

	testApp.Put("/tasks/bulk-action", func(c *fiber.Ctx) error {
		return UpdateBulkTaskHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.name, func(b *testing.B) {
			request := httptest.NewRequest("PUT", "http://localhost.com/tasks/bulk-action", nil)
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				response, _ := testApp.Test(request)
				compareStatusCodes(b, response, scenario)
			}
			b.StopTimer()
		})
	}
}

func compareStatusCodes(b *testing.B, response *http.Response, scenario scenario) {
	if response.StatusCode != scenario.statusCode {
		b.Fatalf("Expected status code: %d, instead got: %d", scenario.statusCode, response.StatusCode)
	}
}
