package model

import "time"

// Quote is a quote submitted by a user to the application
type Quote struct {
	ID          string
	SubmitterID string
	Quotee      string
	Context     string
	Quote       string
	Created     time.Time
}

// Editable returns true if the quote can be edited by the specified user
// Admin users can edit any quote, while non admins can only edit their own quotes within 1 hour of submission
func (q *Quote) Editable(u User) bool {
	return u.Admin || (u.ID == q.SubmitterID && time.Since(q.Created) < time.Hour)
}
