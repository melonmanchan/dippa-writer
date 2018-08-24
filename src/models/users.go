package models

import "github.com/pkg/errors"

// User ...
type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c Client) CreateUserByName(user *User) (int64, error) {
	res, err := c.DB.NamedExec(`
	INSERT INTO users (name) VALUES (:name) ON CONFLICT ON CONSTRAINT constraint_name DO UPDATE SET name = :name RETURNING id;
	`, user)

	if err != nil {
		return -1, errors.Wrap(err, "Could not insert a new user")
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return -1, errors.Wrap(err, "Could not retrieve last insert id")
	}

	return lastId, nil
}
