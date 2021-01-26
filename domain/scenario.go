package domain

import "github.com/DATA-DOG/go-sqlmock"

type Scenario struct {
	StatusCode    int
	RowsAffected  int64
	InsertId      int64
	ScenarioErr   error
	Name          string
	Id            string
	Data          []byte
	ExpectedTasks []Task
	ExpectedSQL   string
	Page          int64
	PerPage       int64
	Rows          *sqlmock.Rows
	Task          Task
	Url           string
}
