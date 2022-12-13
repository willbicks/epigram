package model

import "time"

// MaxQuizAttempts is the maximum number of times a user can submit an entry quiz.
// If this number is exceeded, the user will be blocked from continuing.
const MaxQuizAttempts = 5

// User is a user of the application
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

// IsAuthorized returns true if the user is authorized to access the application (they have passed the quiz and are not banned, or they are an admin)
func (u User) IsAuthorized() bool {
	return (u.QuizPassed && !u.Banned) || u.Admin
}

// IsAdmin returns true if the user is an admin
func (u User) IsAdmin() bool {
	return u.Admin
}
