package model

import "time"

const MaxQuizAttempts = 5

type User struct {
	ID         string
	Name       string
	Email      string
	PictureURL string
	Created    time.Time
	QuizPassed bool
	// QuizAttempts represents the number of times the user has submitted an entry quiz.
	QuizAttempts int8
	Banned       bool
	Admin        bool
}

func (u User) IsAuthorized() bool {
	return (u.QuizPassed && !u.Banned) || u.Admin
}

func (u User) IsAdmin() bool {
	return u.Admin
}
