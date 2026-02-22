package users

import (
	"fmt"
	"practice3/internal/repository/_postgres"
	"practice3/pkg/modules"
)

type Repository struct {
	db *_postgres.Dialect
}

func NewUserRepository(db *_postgres.Dialect) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUsers() ([]modules.User, error) {
	var users []modules.User
	err := r.db.DB.Select(&users, "SELECT id,name,email,age FROM users")
	return users, err
}

func (r *Repository) GetUserByID(id int) (*modules.User, error) {
	var user modules.User
	err := r.db.DB.Get(&user, "SELECT id,name,email,age FROM users WHERE id=$1", id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}

func (r *Repository) CreateUser(user modules.User) (int, error) {
	var id int
	err := r.db.DB.QueryRow(
		"INSERT INTO users(name,email,age) VALUES($1,$2,$3) RETURNING id",
		user.Name, user.Email, user.Age,
	).Scan(&id)
	return id, err
}

func (r *Repository) UpdateUser(id int, user modules.User) error {
	res, err := r.db.DB.Exec(
		"UPDATE users SET name=$1,email=$2,age=$3 WHERE id=$4",
		user.Name, user.Email, user.Age, id,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *Repository) DeleteUser(id int) error {
	res, err := r.db.DB.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
