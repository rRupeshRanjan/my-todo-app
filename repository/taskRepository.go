package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"my-todo-app/config"
	"my-todo-app/domain"
)

var (
	database *sql.DB
	sqlDriver string
	databaseName string
	logger *zap.Logger
)

const (
	getByIdQuery = "SELECT * FROM tasks WHERE id=?"
	getAllQuery = "SELECT * FROM tasks"
	createQuery = "INSERT INTO tasks (title, description, addedOn, dueBy, status) VALUES (?,?,?,?,?)"
	updateQuery = "UPDATE tasks SET title=?, description=?, addedOn=?, dueBy=?, status=? WHERE id=?"
	deleteQuery = "DELETE FROM tasks WHERE id=?"
	initDbQuery = `CREATE TABLE IF NOT EXISTS tasks (
						id INTEGER PRIMARY KEY, 
						title TEXT, 
						description TEXT, 
						addedOn INTEGER, 
						dueBy INTEGER, 
						status TEXT);`
)

func init() {
	sqlDriver = config.SqlDriver
	databaseName = config.DatabaseName
	logger = config.AppLogger
	setupDatabase()
}

func setupDatabase() {
	database, _ = sql.Open(sqlDriver, databaseName)
	_, err := database.Exec(initDbQuery)
	if err != nil {
		logger.Panic("Failure while initializing database, {}" + err.Error())
	}
}

func GetTaskById(id string) ([]domain.Task, error) {
	rows, err := database.Query(getByIdQuery, id)
	var tasks []domain.Task

	for err == nil && rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.AddedOn, &task.DueBy, &task.Status)
		if err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, err
}

func GetAllTasks() ([]domain.Task, error) {
	rows, err := database.Query(getAllQuery)
	var tasks []domain.Task

	for err == nil && rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.AddedOn, &task.DueBy, &task.Status)
		if err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, err
}

func CreateTask(task domain.Task) (int64, error) {
	statement, err := database.Prepare(createQuery)
	if err == nil {
		result, execError := statement.Exec(task.Title, task.Description, task.AddedOn, task.DueBy, task.Status)
		err = execError
		if err == nil && result != nil {
			return result.LastInsertId()
		}
	}
	return -1, err
}

func UpdateTask(task domain.Task, id string) error {
	statement, err := database.Prepare(updateQuery)
	if err == nil {
		_, err = statement.Exec(task.Title, task.Description, task.AddedOn, task.DueBy, task.Status, id)
	}
	return err
}

func DeleteTask(id string) (int64, error) {
	statement, err := database.Prepare(deleteQuery)
	if err == nil {
		result, ExecError := statement.Exec(id)
		err = ExecError
		if err == nil && result != nil {
			return result.RowsAffected()
		}
	}
	return 0, err
}


// TODO:: Implement this
func SearchTasks(params map[string]string) ([]domain.Task, error) {
	return []domain.Task{}, nil
}
