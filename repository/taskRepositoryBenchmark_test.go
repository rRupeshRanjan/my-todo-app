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
				testUtils.GetRepositoryMocks(testUtils.GetTaskByIdKey, mock, expectedSQL, id, scenario)

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
				testUtils.GetRepositoryMocks(testUtils.GetAllTasksKey, mock, expectedSQL, "", scenario)

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
				testUtils.GetRepositoryMocks(testUtils.CreateTaskKey, mock, expectedSQL, "", scenario)

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
				testUtils.GetRepositoryMocks(testUtils.UpdateTaskKey, mock, expectedSQL, scenario.Id, scenario)

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
	id := "8"

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				testUtils.GetRepositoryMocks(testUtils.DeleteTaskKey, mock, expectedSQL, id, scenario)

				_, err := DeleteTask(id)
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
