package models

import "github.com/pkg/errors"

// User ...
type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c Client) CreateUserByName(user *User) (int64, error) {
	var lastId int64

	err := c.DB.QueryRow(`
	INSERT INTO users (name) VALUES ($1) ON CONFLICT ON CONSTRAINT constraint_name DO UPDATE SET name = $1 RETURNING id;
	`, user.Name).Scan(&lastId)

	if err != nil {
		return -1, errors.Wrap(err, "Could not insert a new user")
	}

	return lastId, nil
}
