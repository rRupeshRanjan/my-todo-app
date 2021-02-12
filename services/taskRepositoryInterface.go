package services

import (
	"my-todo-app/domain"
	"my-todo-app/repository"
)

type TaskRepository struct{}

type TaskRepositoryInterface interface {
	getTaskById(id string) ([]domain.Task, error)
	getAllTasks(page int64, perPage int64) ([]domain.Task, error)
	createTask(task domain.Task) (int64, error)
	updateTask(task domain.Task, id string) error
	deleteTask(id string) (bool, error)
	searchTasks(params map[string]string) ([]domain.Task, error)
}

func (t TaskRepository) getTaskById(id string) ([]domain.Task, error) {
	return repository.GetTaskById(id)
}

func (t TaskRepository) getAllTasks(page int64, perPage int64) ([]domain.Task, error) {
	return repository.GetAllTasks(page, perPage)
}

func (t TaskRepository) createTask(task domain.Task) (int64, error) {
	return repository.CreateTask(task)
}

func (t TaskRepository) updateTask(task domain.Task, id string) error {
	return repository.UpdateTask(task, id)
}

func (t TaskRepository) deleteTask(id string) (bool, error) {
	return repository.DeleteTask(id)
}

func (t TaskRepository) searchTasks(params map[string]string) ([]domain.Task, error) {
	return repository.SearchTasks(params)
}
