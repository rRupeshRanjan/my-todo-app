package testUtils

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/domain"
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
				Name: "should get all tasks",
				ExpectedTasks: []domain.Task{
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
				Rows: sqlmock.NewRows(columns).
					AddRow(8, "sample", "sample", 1, 1, "sample").
					AddRow(88, "sample", "sample", 1, 1, "sample"),
			},
			{
				Name:          "should get no tasks",
				ExpectedTasks: []domain.Task{},
				Rows:          sqlmock.NewRows(columns),
			},
			{
				Name:          "should rollback tx for errors",
				ExpectedTasks: []domain.Task{},
				ScenarioErr:   errors.New("error occurred"),
				Rows:          sqlmock.NewRows(columns),
			},
		}
	case CreateTaskKey:
		return []domain.Scenario{
			{
				Name: "should create task with Id 8",
				Task: domain.Task{
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				},
				InsertId: 8,
			},
			{
				Name: "should rollback tx for errors",
				Task: domain.Task{
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
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
					Id:          8,
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				},
			},
			{
				Name: "should rollback tx for errors",
				Task: domain.Task{
					Id:          8,
					AddedOn:     1,
					DueBy:       1,
					Title:       "sample",
					Description: "sample",
					Status:      "sample",
				},
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
