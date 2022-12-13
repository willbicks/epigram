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
