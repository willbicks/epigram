package model

import (
	"testing"
	"time"
)

func TestUserSession_IsExpired(t *testing.T) {

	tests := []struct {
		name string
		us   UserSession
		now  time.Time
		want bool
	}{
		{
			"Empty session",
			UserSession{},
			time.Now(),
			true,
		},
		{
			"New session",
			UserSession{
				Created: time.Now(),
				Expires: time.Now().Add(time.Hour),
			},
			time.Now(),
			false,
		},
		{
			"Old session",
			UserSession{
				Created: time.Now().Add(-3 * time.Hour),
				Expires: time.Now().Add(-1 * time.Hour),
			},
			time.Now(),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.us.IsExpired(tt.now); got != tt.want {
				t.Errorf("UserSession.IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}
