package models

import "github.com/pkg/errors"

// User ...
type Room struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (c Client) CreateRoomByName(room *Room) (int64, error) {
	var lastId int64

	err := c.DB.QueryRow(`
	INSERT INTO rooms (name) VALUES ($1) ON CONFLICT ON CONSTRAINT rooms_name_key DO UPDATE SET name = $1 RETURNING id;
	`, room.Name).Scan(&lastId)

	if err != nil {
		return -1, errors.Wrap(err, "Could not insert a new room")
	}

	return lastId, nil
}
