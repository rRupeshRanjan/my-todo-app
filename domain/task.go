package domain

type Task struct {
	Id int64
	AddedOn int64
	DueBy int64
	Title string
	Description string
	Status string
}
