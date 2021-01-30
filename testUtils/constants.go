package testUtils

const (
	GetTaskByIdKey = "getTaskById"
	GetAllTasksKey = "getAllTasks"
	CreateTaskKey  = "createTask"
	UpdateTaskKey  = "updateTask"
	DeleteTaskKey  = "deleteTask"
	SearchTaskKey  = "searchTask"
)

var columns = []string{"o_id", "o_title", "o_description", "o_addedOn", "o_dueBy", "o_status"}
