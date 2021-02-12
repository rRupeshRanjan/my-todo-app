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
				Id:          "8",
				Rows:        sqlmock.NewRows(columns).AddRow(8, "sample", "sample", 1, 1, "sample"),
				ExpectedSQL: "SELECT * FROM tasks WHERE id = ?",
			},
			{
				Name:          "should get no tasks",
				ExpectedTasks: []domain.Task{},
				Id:            "8",
				Rows:          sqlmock.NewRows(columns),
				ExpectedSQL:   "SELECT * FROM tasks WHERE id = ?",
			},
			{
				Name:          "should rollback tx for errors",
				ExpectedTasks: []domain.Task{},
				ScenarioErr:   errors.New("error occurred"),
				Rows:          sqlmock.NewRows(columns),
				ExpectedSQL:   "SELECT * FROM tasks WHERE id = ?",
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
				ExpectedSQL: "SELECT * FROM tasks LIMIT 5 OFFSET 5",
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
				ExpectedSQL: "SELECT * FROM tasks",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "sample"),
			},
			{
				Name:          "should get no tasks",
				ExpectedTasks: []domain.Task{},
				Page:          1,
				PerPage:       5,
				ExpectedSQL:   "SELECT * FROM tasks LIMIT 5 OFFSET 5",
				Rows:          sqlmock.NewRows(columns),
			},
			{
				Name:          "should rollback tx for errors",
				ExpectedTasks: []domain.Task{},
				Page:          1,
				PerPage:       5,
				ScenarioErr:   errors.New("error occurred"),
				ExpectedSQL:   "SELECT * FROM tasks LIMIT 5 OFFSET 5",
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
				InsertId:    8,
				ExpectedSQL: "INSERT INTO tasks (title,description,addedOn,dueBy,status) VALUES (?,?,?,?,?)",
			},
			{
				Name: "should rollback tx for errors",
				Task: domain.Task{
					AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				InsertId:    -1,
				ScenarioErr: errors.New("error occurred"),
				ExpectedSQL: "INSERT INTO tasks (title,description,addedOn,dueBy,status) VALUES (?,?,?,?,?)",
			},
		}
	case UpdateTaskKey:
		return []domain.Scenario{
			{
				Name: "should update task with 8",
				Task: domain.Task{
					Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				Id:          "8",
				ExpectedSQL: "UPDATE tasks SET title = ?, description = ?, addedOn = ?, dueBy = ?, status = ? WHERE id = ?",
			},
			{
				Name: "should rollback tx for errors",
				Task: domain.Task{
					Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "sample",
				},
				Id:          "8",
				ScenarioErr: errors.New("error occurred"),
				ExpectedSQL: "UPDATE tasks SET title = ?, description = ?, addedOn = ?, dueBy = ?, status = ? WHERE id = ?",
			},
		}
	case DeleteTaskKey:
		return []domain.Scenario{
			{
				Name:         "should delete task by id",
				RowsAffected: true,
				ExpectedSQL:  "DELETE FROM tasks WHERE id = ?",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "sample").
					AddRow(88, "sample", "sample", 1, 1, "sample"),
			},
			{
				Name:         "should not delete task if not present",
				RowsAffected: false,
				ExpectedSQL:  "DELETE FROM tasks WHERE id = ?",
				Rows:         sqlmock.NewRows([]string{}),
			},
			{
				Name:         "should rollback tx for errors",
				ScenarioErr:  errors.New("error occurred"),
				RowsAffected: false,
				ExpectedSQL:  "DELETE FROM tasks WHERE id = ?",
				Rows:         sqlmock.NewRows([]string{}),
			},
		}
	case SearchTaskKey:
		return []domain.Scenario{
			{
				Name: "should get all tasks with id 8",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
				},
				SearchParams: map[string]string{"id": "8"},
				ExpectedSQL:  "SELECT * FROM tasks WHERE id = ? LIMIT 10 OFFSET 0",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "done"),
			},
			{
				Name: "should get all tasks with addedOn before 10",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
					{Id: 9, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
				},
				SearchParams: map[string]string{"addedOnTo": "10"},
				ExpectedSQL:  "SELECT * FROM tasks WHERE addedOn <= ? LIMIT 10 OFFSET 0",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "done").
					AddRow(9, "sample", "sample", 1, 1, "done"),
			},
			{
				Name: "should get all tasks with addedOn after 10",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 11, DueBy: 11, Title: "sample", Description: "sample", Status: "done"},
					{Id: 9, AddedOn: 11, DueBy: 11, Title: "sample", Description: "sample", Status: "done"},
				},
				SearchParams: map[string]string{"addedOnFrom": "10"},
				ExpectedSQL:  "SELECT * FROM tasks WHERE addedOn >= ? LIMIT 10 OFFSET 0",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 11, 11, "done").
					AddRow(9, "sample", "sample", 11, 11, "done"),
			},
			{
				Name: "should get all tasks with dueBy before 10",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
					{Id: 9, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
				},
				SearchParams: map[string]string{"dueByTo": "10"},
				ExpectedSQL:  "SELECT * FROM tasks WHERE dueBy <= ? LIMIT 10 OFFSET 0",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "done").
					AddRow(9, "sample", "sample", 1, 1, "done"),
			},
			{
				Name: "should get all tasks with dueBy after 10",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 11, DueBy: 11, Title: "sample", Description: "sample", Status: "done"},
					{Id: 9, AddedOn: 11, DueBy: 11, Title: "sample", Description: "sample", Status: "done"},
				},
				SearchParams: map[string]string{"dueByFrom": "10"},
				ExpectedSQL:  "SELECT * FROM tasks WHERE dueBy >= ? LIMIT 10 OFFSET 0",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 11, 11, "done").
					AddRow(9, "sample", "sample", 11, 11, "done"),
			},
			{
				Name: "should get all tasks with status done",
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
					{Id: 9, AddedOn: 1, DueBy: 1, Title: "sample", Description: "sample", Status: "done"},
				},
				SearchParams: map[string]string{"status": "done"},
				ExpectedSQL:  "SELECT * FROM tasks WHERE status = ? LIMIT 10 OFFSET 0",
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "done").
					AddRow(9, "sample", "sample", 1, 1, "done"),
			},
			{
				Name:          "should get no tasks",
				ExpectedTasks: []domain.Task{},
				SearchParams:  map[string]string{"status": "unresolved"},
				ExpectedSQL:   "SELECT * FROM tasks WHERE status = ? LIMIT 10 OFFSET 0",
				Rows:          sqlmock.NewRows(columns),
			},
			{
				Name:          "should default page to 0 and perPage to 10 when garbage value provided",
				ExpectedTasks: []domain.Task{},
				SearchParams:  map[string]string{"page": "-1", "perPage": "the simpsons"},
				ExpectedSQL:   "SELECT * FROM tasks LIMIT 10 OFFSET 0",
				Rows:          sqlmock.NewRows(columns),
			},
			{
				Name:          "should rollback tx for errors",
				ExpectedTasks: []domain.Task{},
				SearchParams:  map[string]string{"page": "0"},
				ScenarioErr:   errors.New("error occurred"),
				ExpectedSQL:   "SELECT * FROM tasks LIMIT 10 OFFSET 0",
				Rows:          sqlmock.NewRows(columns),
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
				ExpectedTasks: []domain.Task{
					{Id: 8, AddedOn: 12345, DueBy: 12345, Title: "sample", Description: "sample", Status: "sample"},
				},
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
					{Id: 8, AddedOn: 12345, DueBy: 12345, Title: "sample", Description: "sample", Status: "sample"},
				},
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
					Id: 1, AddedOn: 123, DueBy: 123, Title: "sample", Description: "sample", Status: "sample",
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
				Data:        []byte(`{addedOn": 12345, "due_by": 12345, "title": "sample`),
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
				RowsAffected: true,
			},
			{
				Name:         "should throw 404 delete task if not present in database",
				StatusCode:   http.StatusNotFound,
				ScenarioErr:  nil,
				RowsAffected: false,
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
				Name: "should search tasks successfully",
				ExpectedTasks: []domain.Task{
					{Id: 6, AddedOn: 12345, DueBy: 12345, Title: "sample 6", Description: "sample", Status: "done"},
					{Id: 7, AddedOn: 12345, DueBy: 12345, Title: "sample 7", Description: "sample", Status: "done"},
					{Id: 8, AddedOn: 12345, DueBy: 12345, Title: "sample 8", Description: "sample", Status: "done"},
				},
				Url:        "http://localhost.com/tasks/search?status=done",
				StatusCode: http.StatusOK,
			},
			{
				Name:          "should give zero tasks for search with status code 200",
				ExpectedTasks: []domain.Task{},
				Url:           "http://localhost.com/tasks/search?status=done",
				StatusCode:    http.StatusOK,
			},
			{
				Name:          "search should give 500 for search task for database errors",
				ExpectedTasks: []domain.Task{},
				ScenarioErr:   errors.New("error while fetching Data"),
				Url:           "http://localhost.com/tasks/search?status=done",
				StatusCode:    http.StatusInternalServerError,
			},
		}
	default:
		return []domain.Scenario{}
	}
}
