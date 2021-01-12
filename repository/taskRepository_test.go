package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/testUtils"
	"reflect"
	"testing"
)

var (
	mockDb *sql.DB
	mock   sqlmock.Sqlmock
	err    error
)

func TestInit(t *testing.T) {
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error opening stub database connection: %s", err)
	}
	setDb(mockDb)
}

func TestGetTaskById(t *testing.T) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.GetTaskByIdKey)
	expectedSQL := "SELECT (.+) FROM tasks WHERE id=\\?"

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			id := scenario.Id
			mock.ExpectBegin()
			mock.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(scenario.Rows).WillReturnError(scenario.ScenarioErr)
			if scenario.ScenarioErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			tasks, err := GetTaskById(id)
			if err != scenario.ScenarioErr {
				t.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if !reflect.DeepEqual(scenario.ExpectedTasks, tasks) {
				t.Error("Expected and actual responses are not same")
			}
		})
	}
}

func TestGetAllTasks(t *testing.T) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.GetAllTasksKey)
	expectedSQL := "SELECT (.+) FROM tasks"

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery(expectedSQL).WillReturnRows(scenario.Rows).WillReturnError(scenario.ScenarioErr)
			if scenario.ScenarioErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			tasks, err := GetAllTasks()
			if err != scenario.ScenarioErr {
				t.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if !reflect.DeepEqual(scenario.ExpectedTasks, tasks) {
				t.Error("Expected and actual responses are not same")
			}
		})
	}
}

func TestCreateTask(t *testing.T) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.CreateTaskKey)
	expectedSQL := "INSERT INTO tasks \\(title, description, addedOn, dueBy, status\\) VALUES \\(\\?,\\?,\\?,\\?,\\?\\)"

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(expectedSQL).
				WithArgs(scenario.Task.Title, scenario.Task.Description, scenario.Task.AddedOn,
					scenario.Task.DueBy, scenario.Task.Status).
				WillReturnResult(sqlmock.NewResult(8, 1)).
				WillReturnError(scenario.ScenarioErr)

			if scenario.ScenarioErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			insertId, err := CreateTask(scenario.Task)
			if err != scenario.ScenarioErr {
				t.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if insertId != scenario.InsertId {
				t.Errorf("Expected insertId: %d, Got: %d", scenario.InsertId, insertId)
			}
		})
	}

}

func TestUpdateTask(t *testing.T) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.UpdateTaskKey)
	expectedSQL := "UPDATE tasks SET title=\\?, description=\\?, addedOn=\\?, dueBy=\\?, status=\\? WHERE id=\\?"

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(expectedSQL).
				WithArgs(scenario.Task.Title, scenario.Task.Description, scenario.Task.AddedOn,
					scenario.Task.DueBy, scenario.Task.Status, "8").
				WillReturnResult(sqlmock.NewResult(8, 1)).
				WillReturnError(scenario.ScenarioErr)

			if scenario.ScenarioErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			err := UpdateTask(scenario.Task, "8")
			if err != scenario.ScenarioErr {
				t.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			}
		})
	}

}

func TestDeleteTask(t *testing.T) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.DeleteTaskKey)
	expectedSQL := "DELETE FROM tasks WHERE id=\\?"

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectExec(expectedSQL).WithArgs("8").
				WillReturnResult(sqlmock.NewResult(-1, scenario.RowsAffected)).
				WillReturnError(scenario.ScenarioErr)
			if scenario.ScenarioErr == nil {
				mock.ExpectCommit()
			} else {
				mock.ExpectRollback()
			}

			rowsAffected, err := DeleteTask("8")
			if err != scenario.ScenarioErr {
				t.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if rowsAffected != scenario.RowsAffected {
				t.Errorf("Failure:: Expected %d row to be affected", scenario.RowsAffected)
			}
		})
	}
}

func TestDestroy(t *testing.T) {
	_ = mockDb.Close()
}
