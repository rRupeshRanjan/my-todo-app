package repository

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"my-todo-app/config"
	"my-todo-app/domain"
	"strconv"
)

var (
	db           *sql.DB
	sqlDriver    string
	databaseName string
	logger       *zap.Logger
	columns      = []string{"title", "description", "addedOn", "dueBy", "status"}
)

const (
	initDbQuery = `CREATE TABLE IF NOT EXISTS tasks (
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

	rows, err := sq.Select("*").
		From("tasks").
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		Query()

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

	query := sq.Select("*").From("tasks")
	if page == -1 || perPage == -1 {
		rows, err = query.RunWith(tx).Query()
	} else {
		rows, err = query.Limit(uint64(perPage)).Offset(uint64(page * perPage)).RunWith(tx).Query()
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

	result, err :=
		sq.Insert("tasks").
			Columns(columns...).
			Values(task.Title, task.Description, task.AddedOn, task.DueBy, task.Status).
			RunWith(tx).
			Exec()

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

	_, err = sq.Update("tasks").
		Set("title", task.Title).
		Set("description", task.Description).
		Set("addedOn", task.AddedOn).
		Set("dueBy", task.DueBy).
		Set("status", task.Status).
		Where(sq.Eq{"id": id}).
		RunWith(tx).
		Exec()
	return err
}

func DeleteTask(id string) (bool, error) {
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer func() {
		switch err {
		case nil:
			_ = tx.Commit()
		default:
			_ = tx.Rollback()
		}
	}()

	result, err :=
		sq.Delete("*").
			From("tasks").
			Where(sq.Eq{"id": id}).
			RunWith(tx).
			Query()
	if err == nil && result != nil {
		return result.Next(), nil
	}

	return false, err
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

	rows, err = getSearchQuery(params).RunWith(tx).Query()

	for err == nil && rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.AddedOn, &task.DueBy, &task.Status)
		if err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, err
}

func getSearchQuery(params map[string]string) sq.SelectBuilder {
	query := sq.Select("*").From("tasks")

	page := getPageNumber(params["page"])
	perPage := getPerPage(params["perPage"])

	for key, value := range params {
		switch key {
		case "id", "status":
			query = query.Where(sq.Eq{key: value})
		case "addedOnFrom", "dueByFrom":
			// strip "From" from key, for correct column names
			key = key[:len(key)-4]
			query = query.Where(sq.GtOrEq{key: value})
		case "addedOnTo", "dueByTo":
			// strip "To" from key, for correct column names
			key = key[:len(key)-2]
			query = query.Where(sq.LtOrEq{key: value})
		}
	}

	return query.Limit(uint64(perPage)).Offset(uint64(page * perPage))
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
