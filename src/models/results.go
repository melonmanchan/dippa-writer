//
//CREATE TABLE google_results (
//    id bigserial primary key,
//    created_at TIMESTAMPTZ,
//    detection_confidence DECIMAL(4) NOT NULL,
//    blurred DECIMAL(4) NOT NULL,
//    joy DECIMAL(4) NOT NULL,
//    sorrow DECIMAL(4) NOT NULL,
//    surprise DECIMAL(4) NOT NULL,
//    image bytea NOT NULL,
//    user_id integer REFERENCES users(id) NOT NULL
//);
//
package models

import (
	"time"

	"github.com/pkg/errors"
)

type GoogleResult struct {
	ID                  int64     `json:"id" db:"id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	DetectionConfidence float32   `json:"detection_confidence" db:"detection_confidence"`
	Blurred             float32   `json:"blurred" db:"blurred"`
	Joy                 float32   `json:"joy" db:"joy"`
	Sorrow              float32   `json:"sorrow" db:"sorrow"`
	Surprise            float32   `json:"surprise" db:"surprise"`
	Image               []byte    `json:"image" db:"image"`
	UserID              int64     `json:"user_id" db:"user_id"`
	RoomID              int64     `json:"room_id" db:"room_id"`
}

type Keyword struct {
	ID        int64   `json:"id" db:"id"`
	Contents  string  `json:"contents" db:"contents"`
	Sentiment float32 `json:"sentiment" db:"sentiment"`
	Relevance float32 `json:"relevance" db:"relevance"`
	Sadness   float32 `json:"sadness" db:"sadness"`
	Joy       float32 `json:"joy" db:"joy"`
	Fear      float32 `json:"fear" db:"fear"`
	Disgust   float32 `json:"disgust" db:"disgust"`
	Anger     float32 `json:"anger" db:"anger"`
}

type WatsonResult struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Contents  string    `json:"contents" db:"contents"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoomID    int64     `json:"room_id" db:"room_id"`
	Keywords  []Keyword `json:"keywords" db:"-"`
}

func (c Client) createGoogleResults(result *GoogleResult) error {
	_, err := c.DB.NamedExec(`
		INSERT INTO google_results (detection_confidence, blurred,
		joy, sorrow, surprise, image, user_id, room_id)
		VALUES (:detection_confidence, :blurred,
		:joy, :sorrow, :surprise, :image, :user_id, :room_id)
	`, result)

	if err != nil {
		return errors.Wrap(err, "Could not insert a new google analytics result")
	}

	return nil
}

func (c Client) createWatsonResult(result *WatsonResult) error {
	tx, err := c.DB.Begin()

	if err != nil {
		return errors.Wrap(err, "Could not instantiate transactionf for inserting new watson result")
	}

	_, err = tx.Exec(`
		INSERT INTO watson_results (contents, user_id, room_id)
		VALUES ($1, $2, $3)
	`, result)

	if err != nil {
		return errors.Wrap(err, "Could not insert a new google analytics result")
	}

	if len(result.Keywords) == 0 {
		return tx.Commit()
	}

	stmt, err := tx.Prepare(`INSERT INTO keywords (contents, sentiment, relevance, sadness, joy, fear, disgust, anger)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)

	if err != nil {
		return errors.Wrap(err, "Could not instantiate prepared statement for bulk insert")
	}

	defer stmt.Close()

	for _, k := range result.Keywords {
		_, err = stmt.Exec(k.Contents, k.Sentiment, k.Relevance, k.Sadness, k.Joy, k.Fear, k.Disgust, k.Anger)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "Could not insert new keyword")
		}
	}

	return nil
}
