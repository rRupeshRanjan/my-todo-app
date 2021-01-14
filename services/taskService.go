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

var (
	taskRepository repository.TaskRepositoryInterface
	logger         *zap.Logger
)

func init() {
	taskRepository = repository.TaskRepository{}
	logger = config.AppLogger
}

func GetTaskByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	task, err := taskRepository.GetTaskById(id)
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
	tasks, err := taskRepository.GetAllTasks()
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

	createdId, err := taskRepository.CreateTask(task)
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

	err = taskRepository.UpdateTask(task, id)
	if err == nil {
		return c.JSON(task)
	}

	logger.Error(fmt.Sprintf("Error while updating task with id=%s: %s", id, err))
	return c.SendStatus(http.StatusInternalServerError)
}

func DeleteTaskByIdHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	rowsAffected, err := taskRepository.DeleteTask(id)
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
