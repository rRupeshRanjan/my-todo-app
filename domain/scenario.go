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
	Rows          *sqlmock.Rows
	Task          Task
}
