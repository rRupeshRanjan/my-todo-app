package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/testUtils"
	"testing"
)

func BenchmarkInit(b *testing.B) {
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		b.Errorf("Error opening stub database connection: %s", err)
	}
	setDb(mockDb)
}

func BenchmarkGetTaskById(b *testing.B) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.GetTaskByIdKey)
	expectedSQL := "SELECT (.+) FROM tasks WHERE id=\\?"

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			id := scenario.Id

			b.StartTimer()
			for i := 0; i < b.N; i++ {
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WithArgs(id).WillReturnRows(scenario.Rows).WillReturnError(scenario.ScenarioErr)
				if scenario.ScenarioErr == nil {
					mock.ExpectCommit()
				} else {
					mock.ExpectRollback()
				}

				_, err := GetTaskById(id)
				if err != scenario.ScenarioErr {
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkGetAllTasks(b *testing.B) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.GetAllTasksKey)
	expectedSQL := "SELECT (.+) FROM tasks"

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				mock.ExpectBegin()
				mock.ExpectQuery(expectedSQL).WillReturnRows(scenario.Rows).WillReturnError(scenario.ScenarioErr)
				if scenario.ScenarioErr == nil {
					mock.ExpectCommit()
				} else {
					mock.ExpectRollback()
				}

				_, err := GetAllTasks()
				if err != scenario.ScenarioErr {
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkCreateTask(b *testing.B) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.CreateTaskKey)
	expectedSQL := "INSERT INTO tasks \\(title, description, addedOn, dueBy, status\\) VALUES \\(\\?,\\?,\\?,\\?,\\?\\)"

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
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

				_, err := CreateTask(scenario.Task)
				if err != scenario.ScenarioErr {
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkUpdateTask(b *testing.B) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.UpdateTaskKey)
	expectedSQL := "UPDATE tasks SET title=\\?, description=\\?, addedOn=\\?, dueBy=\\?, status=\\? WHERE id=\\?"

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
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
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkDeleteTask(b *testing.B) {
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.DeleteTaskKey)
	expectedSQL := "DELETE FROM tasks WHERE id=\\?"

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				mock.ExpectBegin()
				mock.ExpectExec(expectedSQL).WithArgs("8").
					WillReturnResult(sqlmock.NewResult(-1, scenario.RowsAffected)).
					WillReturnError(scenario.ScenarioErr)
				if scenario.ScenarioErr == nil {
					mock.ExpectCommit()
				} else {
					mock.ExpectRollback()
				}

				_, err := DeleteTask("8")
				if err != scenario.ScenarioErr {
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkDestroy(b *testing.B) {
	_ = mockDb.Close()
}
