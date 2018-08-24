package models

import (
	"log"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	types "github.com/melonmanchan/dippa-proto/build/go"

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

func GoogleProtoToGoStructs(res types.GoogleFacialRecognition) (GoogleResult, User, Room) {
	timestampAsTime, _ := ptypes.Timestamp(res.CreatedAt)

	googleResult := GoogleResult{
		ID:                  0,
		CreatedAt:           timestampAsTime,
		DetectionConfidence: res.Emotion.DetectionConfidence,
		Blurred:             res.Emotion.Blurred,
		Joy:                 res.Emotion.Joy,
		Sorrow:              res.Emotion.Sorrow,
		Image:               res.Image,
		UserID:              0,
	}

	user := User{
		Name: res.User.Username,
	}

	room := Room{
		Name: res.User.Room,
	}

	return googleResult, user, room
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

func WatsonProtoToGoStructs(res types.WatsonNLP) (WatsonResult, User, Room) {
	timestampAsTime, _ := ptypes.Timestamp(res.CreatedAt)

	var keywords []Keyword
	protoKeywords := res.GetKeywords()

	for _, k := range protoKeywords {
		keywords = append(keywords, Keyword{
			Contents:  k.Contents,
			Sentiment: k.Sentiment,
			Relevance: k.Relevance,
			Sadness:   k.Emotion.Sadness,
			Joy:       k.Emotion.Joy,
			Fear:      k.Emotion.Fear,
			Disgust:   k.Emotion.Disgust,
			Anger:     k.Emotion.Anger,
		})
	}

	watsonResult := WatsonResult{
		ID:        0,
		CreatedAt: timestampAsTime,
		Contents:  res.Contents,
		UserID:    0,
		Keywords:  keywords,
	}

	user := User{
		Name: res.User.Username,
	}

	room := Room{
		Name: res.User.Room,
	}

	return watsonResult, user, room
}

func (c Client) CreateGoogleResults(result *GoogleResult) error {
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

func (c Client) CreateWatsonResult(result *WatsonResult) error {
	var latestId int
	tx, err := c.DB.Begin()

	if err != nil {
		return errors.Wrap(err, "Could not instantiate transaction for inserting new watson result")
	}

	err = tx.QueryRow(`
		INSERT INTO watson_results (contents, user_id, room_id)
		VALUES ($1, $2, $3) RETURNING id;
	`, result.Contents, result.UserID, result.RoomID).Scan(&latestId)

	if err != nil {
		return errors.Wrap(err, "Could not insert a new watson analytics result")
	}

	if len(result.Keywords) == 0 {
		return tx.Commit()
	}

	stmt, err := tx.Prepare(`INSERT INTO keywords (contents, sentiment, relevance, sadness, joy, fear, disgust, anger, watson_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)

	if err != nil {
		return errors.Wrap(err, "Could not instantiate prepared statement for bulk insert")
	}

	defer stmt.Close()

	log.Println(latestId)

	for _, k := range result.Keywords {
		_, err = stmt.Exec(k.Contents, k.Sentiment, k.Relevance, k.Sadness, k.Joy, k.Fear, k.Disgust, k.Anger, latestId)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "Could not insert new keyword")
		}
	}

	return tx.Commit()
}
