package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"my-todo-app/config"
	"my-todo-app/domain"
	"strconv"
	"strings"
)

var (
	db           *sql.DB
	sqlDriver    string
	databaseName string
	logger       *zap.Logger
)

const (
	getByIdQuery            = "SELECT * FROM tasks WHERE id=?"
	getAllByPaginationQuery = "SELECT * FROM tasks LIMIT ? OFFSET ?"
	getAllQuery             = "SELECT * FROM tasks"
	createQuery             = "INSERT INTO tasks (title, description, addedOn, dueBy, status) VALUES (?,?,?,?,?)"
	updateQuery             = "UPDATE tasks SET title=?, description=?, addedOn=?, dueBy=?, status=? WHERE id=?"
	deleteQuery             = "DELETE FROM tasks WHERE id=?"
	initDbQuery             = `CREATE TABLE IF NOT EXISTS tasks (
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
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
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

func GetTaskById(id string) ([]domain.Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
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

func GetAllTasks(page int64, perPage int64) ([]domain.Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	var rows *sql.Rows
	tasks := []domain.Task{}

	if page == -1 || perPage == -1 {
		rows, err = tx.Query(getAllQuery)
	} else {
		rows, err = tx.Query(getAllByPaginationQuery, perPage, page*perPage)
	}

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
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

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
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	_, err = tx.Exec(updateQuery, task.Title, task.Description, task.AddedOn, task.DueBy, task.Status, id)
	return err
}

func DeleteTask(id string) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	result, err := tx.Exec(deleteQuery, id)
	if err == nil && result != nil {
		return result.RowsAffected()
	}

	return 0, err
}

func SearchTasks(params map[string]string) ([]domain.Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	var rows *sql.Rows
	tasks := []domain.Task{}

	query := getSearchQuery(params)
	logger.Info(fmt.Sprintf("Executing query: %s", query))
	rows, err = tx.Query(query)

	for err == nil && rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.AddedOn, &task.DueBy, &task.Status)
		if err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, err
}

func getSearchQuery(params map[string]string) string {
	page := getPageNumber(params["page"])
	perPage := getPerPage(params["perPage"])

	additionalQuery := strings.Builder{}
	for key, value := range params {
		switch key {
		case "id":
			additionalQuery.WriteString(fmt.Sprintf("id = %s AND ", value))
		case "addedOnFrom":
			additionalQuery.WriteString(fmt.Sprintf("addedOn >= %s AND ", value))
		case "addedOnTo":
			additionalQuery.WriteString(fmt.Sprintf("addedOn <= %s AND ", value))
		case "dueByFrom":
			additionalQuery.WriteString(fmt.Sprintf("dueBy >= %s AND ", value))
		case "dueByTo":
			additionalQuery.WriteString(fmt.Sprintf("dueBy <= %s AND ", value))
		case "status":
			additionalQuery.WriteString(fmt.Sprintf("status = \"%s\" AND ", value))
		}
	}

	baseQuery := strings.Builder{}
	baseQuery.WriteString(getAllQuery)

	if len(additionalQuery.String()) > 0 {
		baseQuery.WriteString(" WHERE ")
		baseQuery.WriteString(additionalQuery.String())
	}

	query := baseQuery.String()
	if len(additionalQuery.String()) > 0 {
		query = query[:len(query)-4]
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", perPage, page*perPage)

	return query
}

func getPerPage(perPageString string) int64 {
	perPage, err := strconv.ParseInt(perPageString, 10, 64)
	if err != nil || perPage <= 0 {
		perPage = 10
	}
	return perPage
}

func getPageNumber(pageString string) int64 {
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil || page < 0 {
		page = 0
	}
	return page
}
