package sqlite

const (
	SelectUserByEmail = "SELECT id, email, password FROM user WHERE email = $1"
	InsertNewUser     = "INSERT INTO user(id, email, password) VALUES($1, $2, $3)"

	SelectTasksByAuthor = "SELECT id, title, description, status, deadline FROM task WHERE author_id = $1"
	InsertNewTask       = "INSERT INTO task(id, author_id, title, description, deadline) VALUES($1, $2, $3, $4, $5)"
	UpdateTaskByID      = "UPDATE task SET title = $1, description = $2, status = $3, deadline = $4 WHERE id = $5"
	DeleteTaskByID      = "DELETE FROM task WHERE id = $1 AND author_id = $2" // Ensure the task belongs to the author before deletion
)
