package testUtils

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/domain"
	"net/http"
)

func GetRepositoryTestScenarios(action string) []domain.Scenario {
	switch action {
	case GetTaskByIdKey:
		return []domain.Scenario{
			{
				Name: "should get only one task by id",
				ExpectedTasks: []domain.Task{{
					Id:          8,
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				}},
				Id:   "8",
				Rows: sqlmock.NewRows(columns).AddRow(8, "sample", "sample", 1, 1, "sample"),
			},
			{
				Name:          "should get no tasks",
				ExpectedTasks: []domain.Task{},
				Id:            "8",
				Rows:          sqlmock.NewRows(columns),
			},
			{
				Name:          "should rollback tx for errors",
				ExpectedTasks: []domain.Task{},
				ScenarioErr:   errors.New("error occurred"),
				Rows:          sqlmock.NewRows(columns),
			},
		}
	case GetAllTasksKey:
		return []domain.Scenario{
			{
				Name: "should get all tasks with page number",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample"},
					{Id: 88, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample"},
				},
				Page:        1,
				PerPage:     5,
				ExpectedSQL: "SELECT (.+) FROM tasks LIMIT \\? OFFSET \\?",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "sample").
					AddRow(88, "sample", "sample", 1, 1, "sample"),
			},
			{
				Name: "should get all tasks with -1 page",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample"},
				},
				Page:        -1,
				PerPage:     1,
				ExpectedSQL: "SELECT (.+) FROM tasks ",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "sample"),
			},
			{
				Name:          "should get no tasks",
				ExpectedTasks: []domain.Task{},
				Page:          1,
				PerPage:       5,
				ExpectedSQL:   "SELECT (.+) FROM tasks LIMIT \\? OFFSET \\?",
				Rows:          sqlmock.NewRows(columns),
			},
			{
				Name:          "should rollback tx for errors",
				ExpectedTasks: []domain.Task{},
				Page:          1,
				PerPage:       5,
				ScenarioErr:   errors.New("error occurred"),
				ExpectedSQL:   "SELECT (.+) FROM tasks LIMIT \\? OFFSET \\?",
				Rows:          sqlmock.NewRows(columns),
			},
		}
	case CreateTaskKey:
		return []domain.Scenario{
			{
				Name: "should create task with Id 8",
				Task: domain.Task{
					AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				InsertId: 8,
			},
			{
				Name: "should rollback tx for errors",
				Task: domain.Task{
					AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				InsertId:    -1,
				ScenarioErr: errors.New("error occurred"),
			},
		}
	case UpdateTaskKey:
		return []domain.Scenario{
			{
				Name: "should update task with 8",
				Task: domain.Task{
					Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				Id: "8",
			},
			{
				Name: "should rollback tx for errors",
				Task: domain.Task{
					Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				Id:          "8",
				ScenarioErr: errors.New("error occurred"),
			},
		}
	case DeleteTaskKey:
		return []domain.Scenario{
			{
				Name:         "should delete task by id",
				RowsAffected: 1,
			},
			{
				Name:         "should not delete task if not present",
				RowsAffected: 0,
			},
			{
				Name:         "should rollback tx for errors",
				ScenarioErr:  errors.New("error occurred"),
				RowsAffected: 0,
			},
		}
	default:
		return []domain.Scenario{}
	}
}

func GetServiceTestScenarios(action string) []domain.Scenario {
	switch action {
	case GetTaskByIdKey:
		return []domain.Scenario{
			{
				Name: "should successfully get task by id",
				ExpectedTasks: []domain.Task{{
					Id:          8,
					AddedOn:     123456789,
					DueBy:       123456789,
					Title:       "sample title",
					Description: "sample description",
					Status:      "sample status",
				}},
				StatusCode: http.StatusOK,
			},
			{
				Name:          "should give 404 for get task by id",
				ExpectedTasks: []domain.Task{},
				StatusCode:    http.StatusNotFound,
			},
			{
				Name:          "should give 500 for get task by id for database errors",
				ExpectedTasks: []domain.Task{},
				ScenarioErr:   errors.New("error while fetching Data"),
				StatusCode:    http.StatusInternalServerError,
			},
		}
	case GetAllTasksKey:
		return []domain.Scenario{
			{
				Name: "should successfully get all tasks",
				ExpectedTasks: []domain.Task{
					{
						Id:      8,
						AddedOn: 123456789,
						DueBy:   123456789, Title: "sample title", Description: "sample description", Status: "sample status",
					}},
				Url:        "http://localhost.com/tasks?page=1&perPage=1",
				StatusCode: http.StatusOK,
			},
			{
				Name:          "should give zero tasks with status code 200",
				ExpectedTasks: []domain.Task{},
				Url:           "http://localhost.com/tasks",
				StatusCode:    http.StatusOK,
			},
			{
				Name:          "should give 500 for get task by id for database errors",
				ExpectedTasks: []domain.Task{},
				ScenarioErr:   errors.New("error while fetching Data"),
				Url:           "http://localhost.com/tasks",
				StatusCode:    http.StatusInternalServerError,
			},
		}
	case CreateTaskKey:
		return []domain.Scenario{
			{
				Name: "should successfully create task",
				Task: domain.Task{
					Id:          1,
					AddedOn:     123,
					DueBy:       123,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				},
				Data:        []byte(`{"added_on": 123, "due_by": 123, "title": "sample", "description": "sample", "status": "sample"}`),
				StatusCode:  http.StatusOK,
				ScenarioErr: nil,
			},
			{
				Name: "should create task overriding the id from request body",
				Task: domain.Task{
					Id: 1, AddedOn: 123, DueBy: 123, Title: "sample", Description: "sample", Status: "sample",
				},
				Data:        []byte(`{"id": 8, "added_on": 123, "due_by": 123, "title": "sample", "description": "sample", "status": "sample"}`),
				StatusCode:  http.StatusOK,
				ScenarioErr: nil,
			},
			{
				Name:        "should throw 400 in create task for malformed body",
				Data:        []byte(`{addedOn": 123456789, "due_by": 123456789, "title": "sample`),
				StatusCode:  http.StatusBadRequest,
				ScenarioErr: nil,
			},
			{
				Name:        "should throw 500 in create task for database errors",
				Data:        []byte(`{"added_on": 123, "due_by": 123, "title": "sample", "description": "sample", "status": "sample"}`),
				StatusCode:  http.StatusInternalServerError,
				ScenarioErr: errors.New("error creating task in database"),
			},
		}
	case UpdateTaskKey:
		return []domain.Scenario{
			{
				Name: "should successfully update a task",
				Task: domain.Task{
					Id: 1, AddedOn: 123, DueBy: 123, Title: "sample", Description: "sample", Status: "sample",
				},
				Data:        []byte(`{"id": 1, "added_on": 123, "due_by": 123, "title": "sample", "description": "sample", "status": "sample"}`),
				StatusCode:  http.StatusOK,
				ScenarioErr: nil,
			},
			{
				Name:        "should throw 400 in update task if IDs are different in URL and request body",
				Data:        []byte(`{"id": 8, "added_on": 123, "due_by": 123, "title": "sample", "description": "sample", "status": "sample"}`),
				StatusCode:  http.StatusBadRequest,
				ScenarioErr: nil,
			},
			{
				Name:        "should throw 400 in update task for malformed body",
				Data:        []byte(`{"id": 1, "added_on": 123, "due_by": 123, "title": "sample}`),
				StatusCode:  http.StatusBadRequest,
				ScenarioErr: nil,
			},
			{
				Name:        "should throw 500 in update task database errors",
				Data:        []byte(`{"id": 1, "added_on": 123, "due_by": 123, "title": "sample"}`),
				StatusCode:  http.StatusInternalServerError,
				ScenarioErr: errors.New("error occurred updating task in database"),
			},
		}
	case DeleteTaskKey:
		return []domain.Scenario{
			{
				Name:         "should successfully delete task",
				StatusCode:   http.StatusNoContent,
				ScenarioErr:  nil,
				RowsAffected: 1,
			},
			{
				Name:         "should throw 404 delete task if not present in database",
				StatusCode:   http.StatusNotFound,
				ScenarioErr:  nil,
				RowsAffected: 0,
			},
			{
				Name:        "should throw 500 in delete task for database errors",
				StatusCode:  http.StatusInternalServerError,
				ScenarioErr: errors.New("error deleting record from database"),
			},
		}
	case SearchTaskKey:
		return []domain.Scenario{
			{
				Name:       "should throw 501",
				StatusCode: http.StatusNotImplemented,
			},
		}
	default:
		return []domain.Scenario{}
	}
}
