package config

import (
	"reflect"
	"testing"
)

func TestParseYAML(t *testing.T) {

	tests := []struct {
		name    string
		yaml    string
		want    Application
		wantErr bool
	}{
		{
			name:    "empty",
			yaml:    ``,
			want:    Application{},
			wantErr: false,
		},
		{
			name: "basic",
			yaml: `
baseURL: http://test:12
title: Epigram
description: Where communities quote
trustProxy: true`,
			want: Application{
				BaseURL:     "http://test:12",
				Title:       "Epigram",
				Description: "Where communities quote",
				TrustProxy:  true,
			},
			wantErr: false,
		},
		{
			name: "oidc-provider",
			yaml: `
OIDCProvider: 
  name: google
  issuerURL: https://accounts.google.com
  clientID: o2131iuy3t5-aVcJfw5i3YWYpQYvn3EuPvr7Jwogau.apps.googleusercontent.com
  clientSecret: GOCSPX-rD_5hv2$s6iKeupiUAD@Ff@aSNasX`,
			want: Application{
				OIDCProvider: OIDCProvider{
					Name:         "google",
					IssuerURL:    "https://accounts.google.com",
					ClientID:     "o2131iuy3t5-aVcJfw5i3YWYpQYvn3EuPvr7Jwogau.apps.googleusercontent.com",
					ClientSecret: "GOCSPX-rD_5hv2$s6iKeupiUAD@Ff@aSNasX",
				},
			},
			wantErr: false,
		},
		{
			name:    "repo-error",
			yaml:    `repo: invalid`,
			want:    Application{},
			wantErr: false,
		},
		{
			name: "repo-inmemory",
			yaml: `repo: inmemory`,
			want: Application{
				Repo: Inmemory,
			},
			wantErr: false,
		},
		{
			name: "repo-sqlite",
			yaml: `repo: sqlite`,
			want: Application{
				Repo: SQLite,
			},
			wantErr: false,
		},
		{
			name: "entryquestions",
			yaml: `entryQuestions: 
  - question: Question 1
    answer: ALIGATOR
  - question: Who was the first President?
    answer: George Washington`,
			want: Application{
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseYAML([]byte(tt.yaml))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseYAML() config: \n\tgot  %+v, \n\twant %+v", got, tt.want)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseyAML() unexpected err: got %v", err)
			}
		})
	}
}
