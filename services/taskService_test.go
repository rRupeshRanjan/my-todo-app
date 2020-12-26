package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"io"
	"my-todo-app/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type taskRepositoryMock struct{}

type scenario struct {
	name         string
	task         domain.Task
	tasks        []domain.Task
	data         []byte
	statusCode   int
	err          error
	rowsAffected int64
}

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

	testApp.Get("/task/:id", func(c *fiber.Ctx) error {
		return GetTaskByIdHandler(c)
	})

	request := httptest.NewRequest("GET", "http://localhost.com/task/8", nil)
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {

			taskRepositoryGetByIdMock = func(id string) ([]domain.Task, error) {
				return scenario.tasks, scenario.err
			}

			response, _ := testApp.Test(request)
			if response.StatusCode == 200 {
				compareResponses(t, scenario.statusCode, scenario.tasks[0], response)
			} else {
				compareResponses(t, scenario.statusCode, nil, response)
			}
		})
	}
}

func TestGetAllTasksHandler(t *testing.T) {
	t.Parallel()
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

	request := httptest.NewRequest("GET", "http://localhost.com/tasks", nil)
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {

			taskRepositoryGetAllTasksMock = func() ([]domain.Task, error) {
				return scenario.tasks, scenario.err
			}

			response, _ := testApp.Test(request)
			compareResponses(t, scenario.statusCode, scenario.tasks, response)
		})
	}
}

func TestCreateTaskHandler(t *testing.T) {
	t.Parallel()
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
		t.Run(scenario.name, func(t *testing.T) {

			taskRepositoryCreateTaskMock = func(task domain.Task) (int64, error) {
				return 1, scenario.err
			}

			request := httptest.NewRequest("POST", "/task", bytes.NewBuffer(scenario.data))
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.statusCode, scenario.task, response)
		})
	}
}

func TestUpdateTaskByIdHandler(t *testing.T) {
	t.Parallel()
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
		t.Run(scenario.name, func(t *testing.T) {

			taskRepositoryUpdateTaskMock = func(task domain.Task, id string) error {
				return scenario.err
			}

			request := httptest.NewRequest("PUT", "/task/1", bytes.NewBuffer(scenario.data))
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.statusCode, scenario.task, response)
		})
	}
}

func TestDeleteTaskByIdHandler(t *testing.T) {
	t.Parallel()
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
	testApp.Delete("/task/:id", func(c *fiber.Ctx) error {
		return DeleteTaskByIdHandler(c)
	})

	request := httptest.NewRequest("DELETE", "/task/8", nil)
	for _, scenario := range scenarios {
		taskRepositoryDeleteTaskMock = func(id string) (int64, error) {
			return scenario.rowsAffected, scenario.err
		}

		response, _ := testApp.Test(request)
		compareResponses(t, scenario.statusCode, nil, response)
	}
}

func TestSearchHandler(t *testing.T) {
	t.Parallel()
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
		t.Run(scenario.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/tasks/search", nil)
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.statusCode, nil, response)
		})
	}
}

func TestUpdateBulkTaskHandler(t *testing.T) {
	t.Parallel()
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
		t.Run(scenario.name, func(t *testing.T) {
			request := httptest.NewRequest("PUT", "/tasks/bulk-action", nil)
			response, _ := testApp.Test(request)
			compareResponses(t, scenario.statusCode, nil, response)
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
