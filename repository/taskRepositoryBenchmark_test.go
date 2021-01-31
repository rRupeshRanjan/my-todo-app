package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"my-todo-app/testUtils"
	"testing"
)

func InitialBenchmarkSetup(b *testing.B) {
	mockDb, mock, err = sqlmock.New()
	if err != nil {
		b.Errorf("Error opening stub database connection: %s", err)
	}
	setDb(mockDb)
}

func BenchmarkGetTaskById(b *testing.B) {
	InitialBenchmarkSetup(b)
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
	_ = mockDb.Close()
}

func BenchmarkGetAllTasks(b *testing.B) {
	InitialBenchmarkSetup(b)
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.GetAllTasksKey)

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				testUtils.GetRepositoryMocks(testUtils.GetAllTasksKey, mock, scenario.ExpectedSQL, "", scenario)

				_, err := GetAllTasks(scenario.Page, scenario.PerPage)
				if err != scenario.ScenarioErr {
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
	_ = mockDb.Close()
}

func BenchmarkCreateTask(b *testing.B) {
	InitialBenchmarkSetup(b)
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
	_ = mockDb.Close()
}

func BenchmarkUpdateTask(b *testing.B) {
	InitialBenchmarkSetup(b)
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
	_ = mockDb.Close()
}

func BenchmarkDeleteTask(b *testing.B) {
	InitialBenchmarkSetup(b)
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
	_ = mockDb.Close()
}

func BenchmarkSearchTasks(b *testing.B) {
	InitialBenchmarkSetup(b)
	scenarios := testUtils.GetRepositoryTestScenarios(testUtils.SearchTaskKey)

	for _, scenario := range scenarios {
		b.Run(scenario.Name, func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				testUtils.GetRepositoryMocks(testUtils.SearchTaskKey, mock, scenario.ExpectedSQL, "", scenario)

				_, err := SearchTasks(scenario.SearchParams)
				if err != scenario.ScenarioErr {
					b.Errorf("Expected error: %s, but got: %s", scenario.ScenarioErr, err)
				}
			}
			b.StopTimer()
		})
	}
	_ = mockDb.Close()
}

func BenchmarkGetPageNumber(b *testing.B) {
	scenarios := map[string]int64{
		"1":       1,
		"-1":      0,
		"garbage": 0,
	}

	for input, expected := range scenarios {
		b.Run("Get page number from string", func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				actual := getPageNumber(input)
				if actual != expected {
					b.Errorf("Expected: %d, Actual: %d", expected, actual)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkGetPerPage(b *testing.B) {
	scenarios := map[string]int64{
		"1":       1,
		"-1":      10,
		"garbage": 10,
	}

	for input, expected := range scenarios {
		b.Run("Get per page from string", func(b *testing.B) {
			b.StartTimer()
			for i := 0; i < b.N; i++ {
				actual := getPerPage(input)
				if actual != expected {
					b.Errorf("Expected: %d, Actual: %d", expected, actual)
				}
			}
			b.StopTimer()
		})
	}
}
