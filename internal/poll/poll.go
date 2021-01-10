package poll

import (
	"time"

	"github.com/google/uuid"
)

// Poll represents a unique poll object which is made up of multiple PollOptions
type Poll struct {
	ID          string   `json:"id" bson:"_id"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description"`
	Options     []Option `json:"options" bson:"options"`
	Votes       int64    `json:"votes" bson:"votes"` // votes is the total number of votes cast for a poll
	Created     int64    `json:"created" bson:"created"`
	Modified    int64    `json:"modified" bson:"modified"`
	Expiry      int64    `json:"expiry" bson:"expiry"`
}

// Option represents an option in a poll
type Option struct {
	ID    string `json:"id" bson:"_id"`
	Votes int64  `json:"votes" bson:"votes"`
	Title string `json:"title" bson:"title"`
}

// NewOption is a builder for the Option struct
func NewOption(title string) Option {
	uid := uuid.New().String()
	return Option{ID: uid, Votes: 0, Title: title}
}

// New is a builder for the Poll type
func New(title string, description string, options []Option, expiry int64) Poll {
	uid := uuid.New().String()
	created := time.Now().Unix()
	return Poll{
		ID:          uid,
		Title:       title,
		Description: description,
		Votes:       0,
		Options:     options,
		Created:     created,
		Modified:    created,
		Expiry:      expiry,
	}
}

// Repository is the interface that all database Connections must implement
type Repository interface {
	GetPollByID(id string) (*Poll, error)
	SavePoll(p *Poll) error
	DeletePoll(p *Poll) (bool, error)
	DeletePollByID(id string) (bool, error)
	AddPollVote(pollID string, optionID string) (bool, error)
}
