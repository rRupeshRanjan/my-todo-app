package domain

type Task struct {
	Id          int64  `json:"id"`
	AddedOn     int64  `json:"added_on"`
	DueBy       int64  `json:"due_by"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (t *Task) SetId(id int64) {
	t.Id = id
}

func (t *Task) SetAddedOn(addedOn int64) {
	t.AddedOn = addedOn
}

func (t *Task) SetDueBy(dueBy int64) {
	t.DueBy = dueBy
}

func (t *Task) SetTitle(title string) {
	t.Title = title
}

func (t *Task) SetDescription(description string) {
	t.Description = description
}

func (t *Task) SetStatus(status string) {
	t.Status = status
}

func (t *Task) GetId() int64 {
	return t.Id
}

func (t *Task) GetAddedOn() int64 {
	return t.AddedOn
}

func (t *Task) GetDueBy() int64 {
	return t.DueBy
}

func (t *Task) GetTitle() string {
	return t.Title
}

func (t *Task) GetDescription() string {
	return t.Description
}

func (t *Task) GetStatus() string {
	return t.Status
}
