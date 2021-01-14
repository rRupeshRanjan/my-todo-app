package repository

import (
	"database/sql"
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
	getAllQuery  = "SELECT * FROM tasks LIMIT ? OFFSET ?"
	createQuery  = "INSERT INTO tasks (title, description, addedOn, dueBy, status) VALUES (?,?,?,?,?)"
	updateQuery  = "UPDATE tasks SET title=?, description=?, addedOn=?, dueBy=?, status=? WHERE id=?"
	deleteQuery  = "DELETE FROM tasks WHERE id=?"
	initDbQuery  = `CREATE TABLE IF NOT EXISTS tasks (
						id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, 
						title TEXT NOT NULL, 
						description TEXT NOT NULL, 
						addedOn BIGINT NOT NULL, 
						dueBy BIGINT NOT NULL, 
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
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = db.Exec(initDbQuery)
	if err != nil {
		logger.Panic("Failure while initializing database, {}" + err.Error())
	}
}

func setDb(database *sql.DB) {
	db = database
}

func getTaskById(id string) ([]domain.Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

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

func getAllTasks(page int64, perPage int64) ([]domain.Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	rows, err := tx.Query(getAllQuery, perPage, page*perPage)
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

func createTask(task domain.Task) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(createQuery, task.Title, task.Description, task.AddedOn, task.DueBy, task.Status)
	if err == nil && result != nil {
		return result.LastInsertId()
	}
	return -1, err
}

func updateTask(task domain.Task, id string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(updateQuery, task.Title, task.Description, task.AddedOn, task.DueBy, task.Status, id)
	return err
}

func deleteTask(id string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		switch err {
		case nil:
			tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	result, err := tx.Exec(deleteQuery, id)
	if err == nil && result != nil {
		return result.RowsAffected()
	}

	return 0, err
}

// TODO:: Implement this
func searchTasks(params map[string]string) ([]domain.Task, error) {
	return []domain.Task{}, nil
}
