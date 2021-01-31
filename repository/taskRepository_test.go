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

func InitialSetup(t *testing.T) {
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		t.Errorf("Error opening stub database connection: %s", err)
	}
	setDb(mockDb)
}

func TestGetTaskById(t *testing.T) {
	InitialSetup(t)
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
	_ = mockDb.Close()
}

func TestGetAllTasks(t *testing.T) {
	InitialSetup(t)
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
	_ = mockDb.Close()
}

func TestCreateTask(t *testing.T) {
	InitialSetup(t)
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
	_ = mockDb.Close()
}

func TestUpdateTask(t *testing.T) {
	InitialSetup(t)
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
	_ = mockDb.Close()
}

func TestDeleteTask(t *testing.T) {
	InitialSetup(t)
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
	_ = mockDb.Close()
}

func TestSearchTasks(t *testing.T) {
	InitialSetup(t)
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.SearchTaskKey)

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			testUtils.GetRepositoryMocks(testUtils.SearchTaskKey, mock, scenario.ExpectedSQL, "", scenario)

			tasks, err := SearchTasks(scenario.SearchParams)
			if err != scenario.ScenarioErr {
				t.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
			} else if mock.ExpectationsWereMet() != nil {
				t.Errorf("Expectations were not met: %s", err)
			} else if !reflect.DeepEqual(scenario.ExpectedTasks, tasks) {
				t.Error("Expected and actual responses are not same")
			}
		})
	}
	_ = mockDb.Close()
}

func TestGetPageNumber(t *testing.T) {
	scenarios := map[string]int64{
		"1":       1,
		"-1":      0,
		"garbage": 0,
	}

	for input, expected := range scenarios {
		t.Run("Get page number from string", func(t *testing.T) {
			actual := getPageNumber(input)

			if actual != expected {
				t.Errorf("Expected: %d, Actual: %d", expected, actual)
			}
		})
	}
}

func TestGetPerPage(t *testing.T) {
	scenarios := map[string]int64{
		"1":       1,
		"-1":      10,
		"garbage": 10,
	}

	for input, expected := range scenarios {
		t.Run("Get per page from string", func(t *testing.T) {
			actual := getPerPage(input)

			if actual != expected {
				t.Errorf("Expected: %d, Actual: %d", expected, actual)
			}
		})
	}
}
