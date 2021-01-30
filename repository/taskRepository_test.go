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
			testUtils.GetRepositoryMocks(testUtils.GetTaskByIdKey, mock, expectedSQL, id, scenario)

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

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			testUtils.GetRepositoryMocks(testUtils.GetAllTasksKey, mock, scenario.ExpectedSQL, "", scenario)

			tasks, err := GetAllTasks(scenario.Page, scenario.PerPage)
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
			testUtils.GetRepositoryMocks(testUtils.CreateTaskKey, mock, expectedSQL, "", scenario)

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
			testUtils.GetRepositoryMocks(testUtils.UpdateTaskKey, mock, expectedSQL, scenario.Id, scenario)

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
	id := "8"

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			testUtils.GetRepositoryMocks(testUtils.DeleteTaskKey, mock, expectedSQL, id, scenario)

			rowsAffected, err := DeleteTask(id)
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
