package services

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"my-todo-app/config"
	"my-todo-app/domain"
	"my-todo-app/repository"
	"net/http"
	"strconv"
)

type TaskRepository struct{}

type TaskRepositoryInterface interface {
	getTaskById(id string) ([]domain.Task, error)
	getAllTasks() ([]domain.Task, error)
	createTask(task domain.Task) (int64, error)
	updateTask(task domain.Task, id string) error
	deleteTask(id string) (int64, error)
	searchTasks(params map[string]string) ([]domain.Task, error)
}

func (t TaskRepository) getTaskById(id string) ([]domain.Task, error) {
	return repository.GetTaskById(id)
}

func (t TaskRepository) getAllTasks() ([]domain.Task, error) {
	return repository.GetAllTasks()
}

func (t TaskRepository) createTask(task domain.Task) (int64, error) {
	return repository.CreateTask(task)
}

func (t TaskRepository) updateTask(task domain.Task, id string) error {
	return repository.UpdateTask(task, id)
}

func (t TaskRepository) deleteTask(id string) (int64, error) {
	return repository.DeleteTask(id)
}

func (t TaskRepository) searchTasks(params map[string]string) ([]domain.Task, error) {
	return repository.SearchTasks(params)
}

var (
	taskRepository TaskRepositoryInterface
	logger         *zap.Logger
)

func init() {
	taskRepository = TaskRepository{}
	logger = config.AppLogger
}

func GetTaskByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := taskRepository.getTaskById(id)
	if err == nil {
		if len(task) == 0 {
			logger.Info(fmt.Sprintf("No task found with id: %s", id))
			return c.SendStatus(http.StatusNotFound)
		}
		return c.JSON(task[0])
	}

	logger.Error(fmt.Sprintf("Error fetching task with id=%s: %s", id, err))
	return c.SendStatus(http.StatusInternalServerError)
}

func GetAllTasksHandler(c *fiber.Ctx) error {
	tasks, err := taskRepository.getAllTasks()
	if err == nil {
		logger.Info(fmt.Sprintf("No. of expectedTasks fetched: %d", len(tasks)))
		return c.JSON(tasks)
	}

	logger.Error(fmt.Sprintf("Error fetching all expectedTasks: %s", err))
	return c.SendStatus(http.StatusInternalServerError)
}

func CreateTaskHandler(c *fiber.Ctx) error {
	var task domain.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil {
		logger.Error(fmt.Sprintf("Error converting json to valid task body: %s", err))
		return c.SendStatus(http.StatusBadRequest)
	}

	createdId, err := taskRepository.createTask(task)
	if err == nil {
		task.Id = createdId
		return c.JSON(task)
	}

	logger.Error(fmt.Sprintf("Error creating task: %s", err))
	return c.SendStatus(http.StatusInternalServerError)
}

func UpdateTaskByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	var task domain.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil || strconv.FormatInt(task.Id, 10) != id {
		logger.Error("Bad data passed for update, or id in body is different from id in URL")
		return c.SendStatus(http.StatusBadRequest)
	}

	err = taskRepository.updateTask(task, id)
	if err == nil {
		return c.JSON(task)
	}

	logger.Error(fmt.Sprintf("Error while updating task with id=%s: %s", id, err))
	return c.SendStatus(http.StatusInternalServerError)
}

func DeleteTaskByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	rowsAffected, err := taskRepository.deleteTask(id)
	if err == nil {
		if rowsAffected == 0 {
			logger.Info(fmt.Sprintf("No task found with id: %s for deletion", id))
			return c.SendStatus(http.StatusNotFound)
		}
		logger.Info(fmt.Sprintf("Deleted task with id: %s", id))
		return c.SendStatus(http.StatusNoContent)
	}

	logger.Error(fmt.Sprintf("Error deleting task with id=%s : %s", id, err))
	return c.SendStatus(http.StatusInternalServerError)
}

// TODO: Implement this
func SearchHandler(c *fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).SendString("Not yet implemented")
}

// TODO: Implement this
func UpdateBulkTaskHandler(c *fiber.Ctx) error {
	return c.Status(http.StatusNotImplemented).SendString("Not yet implemented")
}
