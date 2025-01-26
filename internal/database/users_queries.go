package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Password  string
	IsAdm     bool
	IsActive  bool
}

type CreateUserParam struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Password  string
	IsAdm     bool
	IsActive  bool
}

func (q *Queries) CheckDuplicatedEmail(userEmail string) bool {
	query := "SELECT email FROM users WHERE email = $1"

	var email string
	q.db.QueryRow(query, userEmail).Scan(&email)

	return email != ""
}

func (q *Queries) CreateUser(arg CreateUserParam) (User, error) {
	query := "INSERT INTO users(created_at, updated_at, name, email, password, is_adm) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at, name, email, password, is_adm, is_active;"

	row := q.db.QueryRow(query, arg.CreatedAt, arg.UpdatedAt, arg.Name, arg.Email, arg.Password, arg.IsAdm)

	var user User

	if err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.IsAdm, &user.IsActive); err != nil {
		return User{}, err
	}

	return user, nil
}

type ListAllUsersParams struct {
	Offset int
	Limit  int
}

func (q *Queries) ListAllUsersInDB(arg ListAllUsersParams) ([]User, error) {
	query := "SELECT id, created_at, updated_at, name, email, password, is_adm, is_active FROM users OFFSET $1 LIMIT $2;"

	rows, err := q.db.Query(query, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Password, &user.IsAdm, &user.IsActive); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
