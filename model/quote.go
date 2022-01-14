package model

import "time"

type Quote struct {
	ID          string
	SubmitterID string
	Quotee      string
	Context     string
	Quote       string
	Created     time.Time
}
