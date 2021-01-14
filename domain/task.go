package domain

type Task struct {
	Id          int64  `json:"id"`
	AddedOn     int64  `json:"added_on"`
	DueBy       int64  `json:"due_by"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
