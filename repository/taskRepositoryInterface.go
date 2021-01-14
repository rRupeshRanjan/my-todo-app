package repository

import (
	"my-todo-app/domain"
)

type TaskRepository struct{}

type TaskRepositoryInterface interface {
	GetTaskById(id string) ([]domain.Task, error)
	GetAllTasks() ([]domain.Task, error)
	CreateTask(task domain.Task) (int64, error)
	UpdateTask(task domain.Task, id string) error
	DeleteTask(id string) (int64, error)
	SearchTasks(params map[string]string) ([]domain.Task, error)
}

func (t TaskRepository) GetTaskById(id string) ([]domain.Task, error) {
	return getTaskById(id)
}

func (t TaskRepository) GetAllTasks() ([]domain.Task, error) {
	return getAllTasks()
}

func (t TaskRepository) CreateTask(task domain.Task) (int64, error) {
	return createTask(task)
}

func (t TaskRepository) UpdateTask(task domain.Task, id string) error {
	return updateTask(task, id)
}

func (t TaskRepository) DeleteTask(id string) (int64, error) {
	return deleteTask(id)
}

func (t TaskRepository) SearchTasks(params map[string]string) ([]domain.Task, error) {
	return searchTasks(params)
}
