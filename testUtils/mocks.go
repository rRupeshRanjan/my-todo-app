package testUtils

import (
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/domain"
	"strconv"
)

func GetRepositoryMocks(action string, mock sqlmock.Sqlmock, expectedSQL string, id string, scenario domain.Scenario) {
	mock.ExpectBegin()

	integerId, _ := strconv.ParseInt(scenario.Id, 10, 64)
	switch action {
	case GetTaskByIdKey:
		mock.ExpectQuery(expectedSQL).
			WithArgs(id).
			WillReturnRows(scenario.Rows).
			WillReturnError(scenario.ScenarioErr)

	case GetAllTasksKey:
		mock.ExpectQuery(expectedSQL).
			WillReturnRows(scenario.Rows).
			WillReturnError(scenario.ScenarioErr)

	case CreateTaskKey:
		mock.ExpectExec(expectedSQL).
			WithArgs(scenario.Task.Title, scenario.Task.Description, scenario.Task.AddedOn,
				scenario.Task.DueBy, scenario.Task.Status).
			WillReturnResult(sqlmock.NewResult(8, 1)).
			WillReturnError(scenario.ScenarioErr)

	case UpdateTaskKey:
		mock.ExpectExec(expectedSQL).
			WithArgs(scenario.Task.Title, scenario.Task.Description, scenario.Task.AddedOn,
				scenario.Task.DueBy, scenario.Task.Status, scenario.Id).
			WillReturnResult(sqlmock.NewResult(integerId, 1)).
			WillReturnError(scenario.ScenarioErr)

	case DeleteTaskKey:
		mock.ExpectExec(expectedSQL).WithArgs(id).
			WillReturnResult(sqlmock.NewResult(-1, scenario.RowsAffected)).
			WillReturnError(scenario.ScenarioErr)
	}

	if scenario.ScenarioErr == nil {
		mock.ExpectCommit()
	} else {
		mock.ExpectRollback()
	}
}
