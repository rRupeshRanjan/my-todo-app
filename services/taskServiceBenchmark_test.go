package services

import (
	"bytes"
	"github.com/gofiber/fiber/v2"
	"my-todo-app/domain"
	"my-todo-app/testUtils"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func InitialSetup() {
	testApp = fiber.New()
	taskRepository = taskRepositoryMock{}
}

func BenchmarkGetTaskByIdHandler(b *testing.B) {
	InitialSetup()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.GetTaskByIdKey)

	testApp.Get("/task/:id", func(c *fiber.Ctx) error {
		return GetTaskByIdHandler(c)
	})

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
	InitialSetup()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.GetAllTasksKey)

	testApp.Get("/tasks", func(c *fiber.Ctx) error {
		return GetAllTasksHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			taskRepositoryGetAllTasksMock = func(page int64, perPage int64) ([]domain.Task, error) {
				return scenario.ExpectedTasks, scenario.ScenarioErr
			}
			request := httptest.NewRequest("GET", scenario.Url, nil)

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
	InitialSetup()
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
	InitialSetup()
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
	InitialSetup()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.DeleteTaskKey)

	testApp.Delete("/task/:id", func(c *fiber.Ctx) error {
		return DeleteTaskByIdHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			request := httptest.NewRequest("DELETE", "http://localhost.com/task/8", nil)
			taskRepositoryDeleteTaskMock = func(id string) (bool, error) {
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
	InitialSetup()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.SearchTaskKey)

	testApp.Get("/tasks/search", func(c *fiber.Ctx) error {
		return SearchHandler(c)
	})

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			request := httptest.NewRequest("GET", scenario.Url, nil)
			taskRepositorySearchTasksMock = func(params map[string]string) ([]domain.Task, error) {
				return scenario.ExpectedTasks, scenario.ScenarioErr
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

func BenchmarkSingleParamBuildQueryParams(b *testing.B) {
	scenarios := []domain.SearchParamScenario{
		{
			Name: "should add key-value pair to map",
			Key:  "sample key", Value: "sample value", Exists: true,
		},
		{
			Name: "should add id to map",
			Key:  "id", Value: "123", Exists: true,
		},
		{
			Name: "should not add id to map",
			Key:  "id", Value: "", Exists: false,
		},
		{
			Name: "should add status to map",
			Key:  "status", Value: "done", Exists: true,
		},
		{
			Name: "should not add status to map",
			Key:  "status", Value: "", Exists: false,
		},
	}

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				params := map[string]string{}
				buildQueryParams(scenario.Key, scenario.Value, &params)

				setValue, exists := params[scenario.Key]
				if exists != scenario.Exists {
					b.Errorf("Expected key: %s presence in map as: %t, instead got: %t", scenario.Key, scenario.Exists, exists)
				} else if setValue != scenario.Value {
					b.Errorf("Expected value: %s, instead got: %s", scenario.Value, setValue)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkMultiParamBuildQueryParams(b *testing.B) {
	params := map[string]string{}
	entries := []domain.SearchParamScenario{
		{Key: "sample key 1", Value: "sample value 1"},
		{Key: "sample key 2", Value: "sample value 2"},
		{Key: "id", Value: ""},
		{Key: "status", Value: "done"},
	}

	expectedMap := map[string]string{
		"sample key 1": "sample value 1",
		"sample key 2": "sample value 2",
		"status":       "done",
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, entry := range entries {
			buildQueryParams(entry.Key, entry.Value, &params)
		}

		if !reflect.DeepEqual(params, expectedMap) {
			b.Errorf("Expected and computed maps are not equal")
		}
	}
	b.StopTimer()
}

func compareStatusCodes(b *testing.B, response *http.Response, s domain.Scenario) {
	if response.StatusCode != s.StatusCode {
		b.Fatalf("Expected status code: %d, instead got: %d", s.StatusCode, response.StatusCode)
	}
}
