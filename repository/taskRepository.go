package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"my-todo-app/config"
	"my-todo-app/domain"
)

var (
	db           *sql.DB
	sqlDriver    string
	databaseName string
	logger       *zap.Logger
)

const (
	getByIdQuery = "SELECT * FROM tasks WHERE id=?"
	getAllQuery  = "SELECT * FROM tasks"
	createQuery  = "INSERT INTO tasks (title, description, addedOn, dueBy, status) VALUES (?,?,?,?,?)"
	updateQuery  = "UPDATE tasks SET title=?, description=?, addedOn=?, dueBy=?, status=? WHERE id=?"
	deleteQuery  = "DELETE FROM tasks WHERE id=?"
	initDbQuery  = `CREATE TABLE IF NOT EXISTS tasks (
						id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT, 
						title TEXT NOT NULL, 
						description TEXT NOT NULL, 
						addedOn INTEGER NOT NULL, 
						dueBy INTEGER NOT NULL, 
						status TEXT NOT NULL);`
)

func init() {
	sqlDriver = config.SqlDriver
	databaseName = config.DatabaseName
	logger = config.AppLogger

	database := connectDatabase()
	setDb(database)
	initializeTable()
}

func connectDatabase() *sql.DB {
	database, err := sql.Open(sqlDriver, databaseName)
	if err != nil {
		panic(err)
	}

	return database
}

func initializeTable() {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer completeTransaction(err, tx)

	_, err = db.Exec(initDbQuery)
	if err != nil {
		logger.Panic("Failure while initializing database, {}" + err.Error())
	}
}

func setDb(database *sql.DB) {
	db = database
}

func GetTaskById(id string) ([]domain.Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer completeTransaction(err, tx)

	rows, err := tx.Query(getByIdQuery, id)
	tasks := []domain.Task{}

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
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer completeTransaction(err, tx)

	rows, err := tx.Query(getAllQuery)
	tasks := []domain.Task{}

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
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	defer completeTransaction(err, tx)

	result, err := tx.Exec(createQuery, task.Title, task.Description, task.AddedOn, task.DueBy, task.Status)
	if err == nil && result != nil {
		return result.LastInsertId()
	}
	return -1, err
}

func UpdateTask(task domain.Task, id string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer completeTransaction(err, tx)

	_, err = tx.Exec(updateQuery, task.Title, task.Description, task.AddedOn, task.DueBy, task.Status, id)
	return err
}

func DeleteTask(id string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer completeTransaction(err, tx)

	result, err := tx.Exec(deleteQuery, id)
	if err == nil && result != nil {
		return result.RowsAffected()
	}

	return 0, err
}

// TODO:: Implement this
func SearchTasks(params map[string]string) ([]domain.Task, error) {
	return []domain.Task{}, nil
}

func completeTransaction(err error, tx *sql.Tx) {
	switch err {
	case nil:
		err = tx.Commit()
	default:
		err = tx.Rollback()
	}

	if err != nil {
		logger.Error(fmt.Sprintf("Error completing the transaction: %s", err))
	}
}
