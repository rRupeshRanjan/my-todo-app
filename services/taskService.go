package services

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"my-todo-app/config"
	"my-todo-app/domain"
	"my-todo-app/repository"
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
	taskRepository = TaskRepository{}
	logger         = config.AppLogger
)

func GetTaskByIdHandler(c *fiber.Ctx) error {
	setApplicationJsonHeader(c)

	id := c.Params("id")
	task, err := taskRepository.getTaskById(id)
	if err == nil {
		if len(task) == 0 {
			logger.Info(fmt.Sprintf("No task found with id: %s", id))
			return c.Status(404).Send(nil)
		} else {
			var body []byte
			body, err = json.Marshal(task[0])
			if err == nil {
				return c.Status(200).Send(body)
			}
		}
	}

	logger.Error(fmt.Sprintf("Error fetching task with id=%s: %s", id, err))
	return c.Status(500).Send(nil)
}

func GetAllTasksHandler(c *fiber.Ctx) error {
	setApplicationJsonHeader(c)

	tasks, err := taskRepository.getAllTasks()
	if err == nil {
		var body []byte
		body, err = json.Marshal(tasks)
		if err == nil {
			logger.Info(fmt.Sprintf("No. of tasks fetched: %d", len(tasks)))
			return c.Status(200).Send(body)
		}
	}

	logger.Error(fmt.Sprintf("Error fetching all tasks: %s", err))
	return c.Status(500).Send(nil)
}

func CreateTaskHandler(c *fiber.Ctx) error {
	c.Set("content-type", "application/json")

	var task domain.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil {
		logger.Error(fmt.Sprintf("Error converting json to valid task body: %s", err))
		return c.Status(400).Send(nil)
	}

	createdId, err := taskRepository.createTask(task)
	if err == nil {
		var body []byte
		task.Id = createdId
		body, err = json.Marshal(task)
		if err == nil {
			logger.Info(fmt.Sprintf("Created task with id: %d", createdId))
			return c.Status(200).Send(body)
		}
	}

	logger.Error(fmt.Sprintf("Error creating task: %s", err))
	return c.Status(500).Send(nil)
}

func UpdateTaskByIdHandler(c *fiber.Ctx) error {
	setApplicationJsonHeader(c)
	id := c.Params("id")

	var task domain.Task
	err := json.Unmarshal(c.Body(), &task)
	if err != nil || strconv.FormatInt(task.Id, 10) != id {
		logger.Error("Bad data passed for update, or id in body is different from id in URL")
		return c.Status(400).Send(nil)
	}

	err = taskRepository.updateTask(task, id)
	if err == nil {
		var body []byte
		body, err = json.Marshal(task)
		if err == nil {
			logger.Info(fmt.Sprintf("Updated task with id: %d", id))
			return c.Status(200).Send(body)
		}
	}

	logger.Error(fmt.Sprintf("Error while updating task with id=%s: %s", id, err))
	return c.Status(500).Send(nil)
}

func DeleteTaskByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	rowsAffected, err := taskRepository.deleteTask(id)
	if err == nil {
		if rowsAffected == 0 {
			logger.Info(fmt.Sprintf("No task found with id: %s for deletion", id))
			return c.Status(404).Send(nil)
		}
		logger.Info(fmt.Sprintf("Deleted task with id: %s", id))
		return c.Status(204).Send(nil)
	}

	logger.Error(fmt.Sprintf("Error deleting task with id=%s : %s", id, err))
	return c.Status(500).Send(nil)
}

// TODO: Implement this
func SearchHandler(c *fiber.Ctx) error {
	return c.Status(501).SendString("Not yet implemented")
}

// TODO: Implement this
func UpdateBulkTaskHandler(c *fiber.Ctx) error {
	return c.Status(501).SendString("Not yet implemented")
}

func setApplicationJsonHeader(c *fiber.Ctx) {
	c.Set("content-type", "application/json")
}
