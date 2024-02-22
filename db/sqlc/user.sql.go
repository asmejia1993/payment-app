package db

import "context"

const createUser = `-- name: CreateUser :one
INSERT INTO users (
    first_name,
	last_name,
    email,
	username,
	user_type,
    hashed_password
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, first_name, last_name, email, username, user_type, hashed_password, password_changed_at, created_at
`
const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, email, username, user_type, hashed_password, password_changed_at, created_at FROM users
WHERE email = $1 LIMIT 1
`

type CreateUserParams struct {
	FistName       string `json:"fistName"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
	Username       string `json:"username"`
	UserType       string `json:"userType"`
	HashedPassword string `json:"hashedPassword"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"Password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.FistName,
		arg.LastName,
		arg.Email,
		arg.Username,
		arg.UserType,
		arg.HashedPassword,
	)

	var u User
	err := row.Scan(
		&u.Id,
		&u.FistName,
		&u.LastName,
		&u.Email,
		&u.Username,
		&u.UserType,
		&u.HashedPassword,
		&u.PasswordChangedAt,
		&u.CreatedAt,
	)
	return u, err
}

func (q *Queries) GetUser(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, email)

	var u User
	err := row.Scan(
		&u.Id,
		&u.FistName,
		&u.LastName,
		&u.Email,
		&u.Username,
		&u.UserType,
		&u.HashedPassword,
		&u.PasswordChangedAt,
		&u.CreatedAt,
	)
	return u, err
}
