package models

import "github.com/pkg/errors"

// User ...
type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c Client) createUserByName(user *User) error {
	_, err := c.DB.NamedExec(`
	INSERT INTO users (name) VALUES (:name);
	`, user)

	if err != nil {
		return errors.Wrap(err, "Could not insert a new user")
	}

	return nil
}

// Query(`INSERT INTO mytable (field1, field2) VALUES (:f1, :f2), (:f2, :f2)`, things)
