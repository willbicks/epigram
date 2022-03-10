// package config is responsible for processing global, application configuration specified at the time of execution.

// Supports both reading from .yml files, and environment variables, the latter taking priority over the former.

package config

import (
	"reflect"
	"testing"
)

func TestApplication_merge(t *testing.T) {
	tests := []struct {
		name  string
		base  Application
		layer Application
		want  Application
	}{
		{
			name:  "default-blank",
			base:  Default,
			layer: Application{},
			want:  Default,
		},
		{
			name: "default-partial_overwrite",
			base: Default,
			layer: Application{
				BaseURL:     "http://test",
				Title:       "Epigram",
				Description: "Epigram example",
				OIDCProvider: OIDCProvider{
					Name:         "test",
					IssuerURL:    "https://accounts.google.com",
					ClientID:     "ididid",
					ClientSecret: "secretsecret",
				},
				EntryQuestions: []EntryQuestion{
					{
						Question: "Question 1",
						Answer:   "ALIGATOR",
					},
					{
						Question: "Who was the first President?",
						Answer:   "George Washington",
					},
				},
			},
			want: Application{
				Address:     Default.Address,
				Port:        Default.Port,
				BaseURL:     "http://test",
				Title:       "Epigram",
				Description: "Epigram example",
				Repo:        Default.Repo,
				DBLoc:       Default.DBLoc,
				TrustProxy:  Default.TrustProxy,
				OIDCProvider: OIDCProvider{
					Name:         "test",
					IssuerURL:    "https://accounts.google.com",
					ClientID:     "ididid",
					ClientSecret: "secretsecret",
				},
				EntryQuestions: []EntryQuestion{
					{
						Question: "Question 1",
						Answer:   "ALIGATOR",
					},
					{
						Question: "Who was the first President?",
						Answer:   "George Washington",
					},
				},
			},
		},
		{
			name: "default-complete_overwrite",
			base: Default,
			layer: Application{
				Address:     "1.2.3.4",
				Port:        303,
				BaseURL:     "http://test",
				Title:       "Epigram",
				Description: "Epigram example",
				Repo:        SQLite,
				DBLoc:       "/var/rando",
				TrustProxy:  true,
				OIDCProvider: OIDCProvider{
					Name:         "test",
					IssuerURL:    "https://accounts.google.com",
					ClientID:     "ididid",
					ClientSecret: "secretsecret",
				},
				EntryQuestions: []EntryQuestion{
					{
						Question: "Question 1",
						Answer:   "ALIGATOR",
					},
					{
						Question: "Who was the first President?",
						Answer:   "George Washington",
					},
				},
			},
			want: Application{
				Address:     "1.2.3.4",
				Port:        303,
				BaseURL:     "http://test",
				Title:       "Epigram",
				Description: "Epigram example",
				Repo:        SQLite,
				DBLoc:       "/var/rando",
				TrustProxy:  true,
				OIDCProvider: OIDCProvider{
					Name:         "test",
					IssuerURL:    "https://accounts.google.com",
					ClientID:     "ididid",
					ClientSecret: "secretsecret",
				},
				EntryQuestions: []EntryQuestion{
					{
						Question: "Question 1",
						Answer:   "ALIGATOR",
					},
					{
						Question: "Who was the first President?",
						Answer:   "George Washington",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.base.merge(tt.layer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Application.merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
