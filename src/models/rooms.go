package models

import "github.com/pkg/errors"

// User ...
type Room struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c Client) createRoomByName(room *Room) error {
	_, err := c.DB.NamedExec(`
	INSERT INTO rooms (name) VALUES (:name);
	`, room)

	if err != nil {
		return errors.Wrap(err, "Could not insert a new room")
	}

	return nil
}
