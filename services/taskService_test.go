package services

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"io"
	"my-todo-app/domain"
	"my-todo-app/testUtils"
	"net/http"
	"net/http/httptest"
	"testing"
)

type taskRepositoryMock struct{}

var (
	taskRepositoryGetByIdMock     func(id string) ([]domain.Task, error)
	taskRepositoryGetAllTasksMock func() ([]domain.Task, error)
	taskRepositoryCreateTaskMock  func(task domain.Task) (int64, error)
	taskRepositoryUpdateTaskMock  func(task domain.Task, id string) error
	taskRepositoryDeleteTaskMock  func(id string) (int64, error)
	taskRepositorySearchTasksMock func(params map[string]string) ([]domain.Task, error)

	testApp = fiber.New()
)

func (t taskRepositoryMock) getTaskById(id string) ([]domain.Task, error) {
	return taskRepositoryGetByIdMock(id)
}

func (t taskRepositoryMock) getAllTasks() ([]domain.Task, error) {
	return taskRepositoryGetAllTasksMock()
}

func (t taskRepositoryMock) createTask(task domain.Task) (int64, error) {
	return taskRepositoryCreateTaskMock(task)
}

func (t taskRepositoryMock) updateTask(task domain.Task, id string) error {
	return taskRepositoryUpdateTaskMock(task, id)
}

func (t taskRepositoryMock) deleteTask(id string) (int64, error) {
	return taskRepositoryDeleteTaskMock(id)
}

func (t taskRepositoryMock) searchTasks(params map[string]string) ([]domain.Task, error) {
	return taskRepositorySearchTasksMock(params)
}

func TestSetup(t *testing.T) {
	taskRepository = taskRepositoryMock{}
}

func TestGetTaskByIdHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.GetTaskByIdKey)

	testApp.Get("/task/:id", func(c *fiber.Ctx) error {
		return GetTaskByIdHandler(c)
	})

	request := httptest.NewRequest("GET", "http://localhost.com/task/8", nil)
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {

			taskRepositoryGetByIdMock = func(id string) ([]domain.Task, error) {
				return scenario.ExpectedTasks, scenario.ScenarioErr
			}

			response, _ := testApp.Test(request)
			if response.StatusCode == 200 {
				compareResponses(t, scenario.StatusCode, scenario.ExpectedTasks[0], response)
			} else {
				compareResponses(t, scenario.StatusCode, nil, response)
			}
		})
	}
}

func TestGetAllTasksHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.GetAllTasksKey)

	testApp.Get("/expectedTasks", func(c *fiber.Ctx) error {
		return GetAllTasksHandler(c)
	})

	request := httptest.NewRequest("GET", "http://localhost.com/expectedTasks", nil)
	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {

			taskRepositoryGetAllTasksMock = func() ([]domain.Task, error) {
				return scenario.ExpectedTasks, scenario.ScenarioErr
			}

			response, _ := testApp.Test(request)
			compareResponses(t, scenario.StatusCode, scenario.ExpectedTasks, response)
		})
	}
}

func TestCreateTaskHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.CreateTaskKey)
	testApp.Post("/task", func(c *fiber.Ctx) error {
		return CreateTaskHandler(c)
	})

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {

			taskRepositoryCreateTaskMock = func(task domain.Task) (int64, error) {
				return 1, scenario.ScenarioErr
			}

			request := httptest.NewRequest("POST", "http://localhost.com/task", bytes.NewBuffer(scenario.Data))
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.StatusCode, scenario.Task, response)
		})
	}
}

func TestUpdateTaskByIdHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.UpdateTaskKey)

	testApp.Put("/task/:id", func(c *fiber.Ctx) error {
		return UpdateTaskByIdHandler(c)
	})

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {

			taskRepositoryUpdateTaskMock = func(task domain.Task, id string) error {
				return scenario.ScenarioErr
			}

			request := httptest.NewRequest("PUT", "http://localhost.com/task/1", bytes.NewBuffer(scenario.Data))
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.StatusCode, scenario.Task, response)
		})
	}
}

func TestDeleteTaskByIdHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.DeleteTaskKey)
	testApp.Delete("/task/:id", func(c *fiber.Ctx) error {
		return DeleteTaskByIdHandler(c)
	})

	request := httptest.NewRequest("DELETE", "http://localhost.com/task/8", nil)
	for _, scenario := range scenarios {
		taskRepositoryDeleteTaskMock = func(id string) (int64, error) {
			return scenario.RowsAffected, scenario.ScenarioErr
		}

		response, _ := testApp.Test(request)
		compareResponses(t, scenario.StatusCode, nil, response)
	}
}

func TestSearchHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.SearchTaskKey)

	testApp.Get("/expectedTasks/search", func(c *fiber.Ctx) error {
		return SearchHandler(c)
	})

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "http://localhost.com/expectedTasks/search", nil)
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.StatusCode, nil, response)
		})
	}
}

func TestUpdateBulkTaskHandler(t *testing.T) {
	t.Parallel()
	scenarios := testUtils.GetServiceTestScenarios(testUtils.BulkUpdateTaskKey)

	testApp.Put("/expectedTasks/bulk-action", func(c *fiber.Ctx) error {
		return UpdateBulkTaskHandler(c)
	})

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			request := httptest.NewRequest("PUT", "http://localhost.com/expectedTasks/bulk-action", nil)
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.StatusCode, nil, response)
		})
	}
}

func compareResponses(t *testing.T, expectedStatusCode int, expectedBody interface{}, response *http.Response) {
	if response.StatusCode != expectedStatusCode {
		t.Errorf("Expected status code: %d, Got: %d", expectedStatusCode, response.StatusCode)
	}

	if response.StatusCode == http.StatusOK {
		actual := getStringFromResponseBody(response.Body)
		expected := getStringFromStruct(expectedBody)
		if actual != expected {
			logMisMatchedData(t, expected, actual)
		}
	}
}

func logMisMatchedData(t *testing.T, expected string, actual string) {
	t.Errorf("\nExpected: %v,\nGot     : %v", expected, actual)
}

func getStringFromResponseBody(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(body)
	return buf.String()
}

func getStringFromStruct(data interface{}) string {
	byteData, _ := json.Marshal(data)
	return string(byteData)
}
