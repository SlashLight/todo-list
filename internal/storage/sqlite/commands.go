package sqlite

const (
	SelectUserByEmail = "SELECT id, email, password FROM user WHERE email = $1"
	InsertNewUser     = "INSERT INTO user(id, email, password) VALUES($1, $2, $3)"
)
