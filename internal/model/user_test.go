package model

import (
	"testing"
)

func TestUser_IsAuthorized(t *testing.T) {
	tests := []struct {
		u    User
		want bool
	}{
		{
			User{
				Name:  "Nick",
				Admin: true,
			},
			true,
		},
		{
			User{
				Name:       "Krystian",
				QuizPassed: true,
			},
			true,
		},
		{
			User{
				Name:       "Outsider",
				QuizPassed: false,
			},
			false,
		},
		{
			User{
				Name:       "DVD",
				QuizPassed: true,
				Banned:     true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.u.Name, func(t *testing.T) {
			if got := tt.u.IsAuthorized(); got != tt.want {
				t.Errorf("User.IsAuthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		u    User
		want bool
	}{
		{
			User{
				Name:  "Nick",
				Admin: true,
			},
			true,
		},
		{
			User{
				Name:       "Krystian",
				QuizPassed: true,
			},
			false,
		},
		{
			User{
				Name:       "Outsider",
				QuizPassed: false,
			},
			false,
		},
		{
			User{
				Name:       "DVD",
				QuizPassed: true,
				Banned:     true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.u.Name, func(t *testing.T) {
			if got := tt.u.IsAdmin(); got != tt.want {
				t.Errorf("User.IsAdmin() = %v, want %v", got, tt.want)
			}
		})
	}
}
