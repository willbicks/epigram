package model

import "time"

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

func (u User) isAuthorized() bool {
	return (u.QuizPassed && !u.Banned) || u.Admin
}

func (u User) isAdmin() bool {
	return u.Admin
}
