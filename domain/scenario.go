package domain

import "github.com/DATA-DOG/go-sqlmock"

type Scenario struct {
	Id            string
	Url           string
	Name          string
	Data          []byte
	Page          int64
	Rows          *sqlmock.Rows
	Task          Task
	SearchParams  map[string]string
	PerPage       int64
	InsertId      int64
	StatusCode    int
	ScenarioErr   error
	RowsAffected  int64
	ExpectedSQL   string
	ExpectedTasks []Task
}

type SearchParamScenario struct {
	Key    string
	Name   string
	Value  string
	Exists bool
}
