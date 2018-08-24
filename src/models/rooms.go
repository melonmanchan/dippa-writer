package models

import "github.com/pkg/errors"

// User ...
type Room struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c Client) CreateRoomByName(room *Room) (int64, error) {
	res, err := c.DB.NamedExec(`
	INSERT INTO rooms (name) VALUES (:name) ON CONFLICT ON CONSTRAINT constraint_name DO UPDATE SET name = :name RETURNING id;
	`, room)

	if err != nil {
		return -1, errors.Wrap(err, "Could not insert a new room")
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return -1, errors.Wrap(err, "Could not retrieve last insert id")
	}

	return lastId, nil
}
