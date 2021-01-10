package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/domain"
	"reflect"
	"testing"
)

type scenario struct {
	name          string
	expectedTasks []domain.Task
	id            string
	rows          *sqlmock.Rows
	scenarioErr   error
	rowsAffected  int64
	task          domain.Task
	insertId      int64
}

var (
	columns = []string{"o_id", "o_title", "o_description", "o_addedOn", "o_dueBy", "o_status"}
	mockDb  *sql.DB
	mock    sqlmock.Sqlmock
	err     error
)

func TestInit(t *testing.T) {
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("An error %s was not expected when opening a stub database connection", err)
	}
	setDb(mockDb)
}

func TestGetTaskById(t *testing.T) {
	scenarios := []scenario{
		{
			name: "should get only one task by id",
			expectedTasks: []domain.Task{{
				Id:          8,
				AddedOn:     1,
				DueBy:       1,
				Title:       "sample",
				Description: "sample",
				Status:      "sample",
			}},
			id:   "8",
			rows: sqlmock.NewRows(columns).AddRow(8, "sample", "sample", 1, 1, "sample"),
		},
		{
			name:          "should get no tasks",
			expectedTasks: []domain.Task{},
			id:            "8",
			rows:          sqlmock.NewRows(columns),
		},
		{
			name:          "should rollback tx for errors",
			expectedTasks: []domain.Task{},
			scenarioErr:   errors.New("error occurred"),
			rows:          sqlmock.NewRows(columns),
		},
	}
	expectedSQL := "SELECT (.+) FROM tasks WHERE id=\\?"

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			id := scenario.id
			mock.ExpectBegin()
			mock.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(scenario.rows)
			mock.ExpectCommit()

			tasks, err := GetTaskById(id)
			if err != nil {
				t.Errorf("Expected no error, but got %s instead while reading from mockDb", err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if !reflect.DeepEqual(scenario.expectedTasks, tasks) {
				t.Error("Expected and actual responses are not same")
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	scenarios := []scenario{
		{
			name: "should get all tasks",
			expectedTasks: []domain.Task{
				{
					Id:          8,
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				},
				{
					Id:          88,
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				},
			},
			rows: sqlmock.NewRows(columns).
				AddRow(8, "sample", "sample", 1, 1, "sample").
				AddRow(88, "sample", "sample", 1, 1, "sample"),
		},
		{
			name:          "should get no tasks",
			expectedTasks: []domain.Task{},
			rows:          sqlmock.NewRows(columns),
		},
	}
	expectedSQL := "SELECT (.+) FROM tasks"

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery(expectedSQL).WillReturnRows(scenario.rows)
			mock.ExpectCommit()

			tasks, err := GetAllTasks()
			if err != nil {
				t.Errorf("Expected no error, but got %s instead while reading from mockDb", err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if !reflect.DeepEqual(scenario.expectedTasks, tasks) {
				t.Error("Expected and actual responses are not same")
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	scenarios := []scenario{
		{
			name: "should create task with Id 8",
			task: domain.Task{
				AddedOn:     1,
				DueBy:       1,
				Title:       "sample",
				Description: "sample",
				Status:      "sample",
			},
			insertId: 8,
		},
	}
	expectedSQL := "INSERT INTO tasks \\(title, description, addedOn, dueBy, status\\) VALUES \\(\\?,\\?,\\?,\\?,\\?\\)"

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(expectedSQL).
				WithArgs(scenario.task.Title, scenario.task.Description, scenario.task.AddedOn,
					scenario.task.DueBy, scenario.task.Status).
				WillReturnResult(sqlmock.NewResult(8, 1)).WillReturnError(scenario.scenarioErr)
			mock.ExpectCommit()

			insertId, err := CreateTask(scenario.task)
			if err != nil {
				t.Errorf("Expected no error, but got %s instead while reading from mockDb", err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if insertId != scenario.insertId {
				t.Errorf("Expected insertId: %d, Got: %d", scenario.insertId, insertId)
			}
		})
	}

}

func TestUpdateTask(t *testing.T) {
	scenarios := []scenario{
		{
			name: "should update task with 8",
			task: domain.Task{
				Id:          8,
				AddedOn:     1,
				DueBy:       1,
				Title:       "sample",
				Description: "sample",
				Status:      "sample",
			},
		},
	}
	expectedSQL := "UPDATE tasks SET title=\\?, description=\\?, addedOn=\\?, dueBy=\\?, status=\\? WHERE id=\\?"

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(expectedSQL).
				WithArgs(scenario.task.Title, scenario.task.Description, scenario.task.AddedOn,
					scenario.task.DueBy, scenario.task.Status, "8").
				WillReturnResult(sqlmock.NewResult(8, 1)).WillReturnError(scenario.scenarioErr)
			mock.ExpectCommit()

			err := UpdateTask(scenario.task, "8")
			if err != nil {
				t.Errorf("Expected no error, but got %s instead while reading from mockDb", err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			}
		})
	}

}

func TestDeleteTask(t *testing.T) {
	scenarios := []scenario{
		{
			name:         "should delete task by id",
			rowsAffected: 1,
		},
		{
			name:         "should not delete task if not present",
			rowsAffected: 0,
		},
	}
	expectedSQL := "DELETE FROM tasks WHERE id=\\?"

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(expectedSQL).WithArgs("8").
				WillReturnResult(sqlmock.NewResult(-1, scenario.rowsAffected))
			mock.ExpectCommit()

			rowsAffected, err := DeleteTask("8")
			if err != nil {
				t.Errorf("Expected no error, but got %s instead while reading from mockDb", err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if rowsAffected != scenario.rowsAffected {
				t.Errorf("Failure:: Expected %d row to be affected", scenario.rowsAffected)
			}
		})
	}
}

func TestDestroy(t *testing.T) {
	_ = mockDb.Close()
}
