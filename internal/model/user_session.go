package model

import "time"

// UserSession describes a session belonging to a user, and allows a cryptographically unique token (ID)
// to be used to identify them. A UserSession also tracks additional information about the user's session,
// such as their IP address, and provides an ability to expire sessions after a given amount of time.
type UserSession struct {
	ID      string
	UserID  string
	Created time.Time
	Expires time.Time
	IP      string
}

// IsExpired returns whether the UserSession has expired relative to the specified current time (now).
func (us UserSession) IsExpired(now time.Time) bool {
	return (us.Expires.IsZero() || us.Expires.Before(now))
}
