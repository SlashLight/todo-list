package sqlite

const (
	CreateUserTable = "CREATE TABLE IF NOT EXISTS user ( id UUID PRIMARY KEY, email TEXT NOT NULL UNIQUE, password TEXT NOT NULL); CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);"

	SelectUserByEmail = "SELECT id, email, password FROM user WHERE email = $1"
	InsertNewUser     = "INSERT INTO user(id, email, password) VALUES($1, $2, $3)"
)
