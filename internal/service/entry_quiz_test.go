package service

import (
	"context"
	"github.com/willbicks/epigram/internal/config"
	"testing"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/model"
)

func TestNewEntryQuizService(t *testing.T) {

	// intSliceContains is a  simple helper to check if a slice contains an element with the provided value
	intSliceContains := func(slice []int, elem int) bool {
		for _, v := range slice {
			if v == elem {
				return true
			}
		}
		return false
	}

	tests := []struct {
		name           string
		entryQuestions []config.EntryQuestion
	}{
		{
			name:           "No questions",
			entryQuestions: []config.EntryQuestion{},
		},
		{
			name: "One question",
			entryQuestions: []config.EntryQuestion{
				{
					Question: "How many chickens can lay an egg?",
					Answer:   "Three",
				},
			},
		},
		{
			name: "Three questions",
			entryQuestions: []config.EntryQuestion{
				{
					Question: "the best place to find a fox",
					Answer:   "wooods",
				},
				{
					Question: "How many chickens can lay an egg?",
					Answer:   "Three",
				},
				{
					Question: "The US's first president",
					Answer:   "washington",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEntryQuizService(tt.entryQuestions)

			// Check that number of questions matches provided
			if len(tt.entryQuestions) != len(got.Questions) {
				t.Errorf("Unexpected number of questions, got %v, want %v", len(got.Questions), len(tt.entryQuestions))
			}

			// Check that IDs are unique
			ids := []int{}
			for _, q := range got.Questions {
				if intSliceContains(ids, q.Id) {
					t.Errorf("Question id %v is non-unique", q.Id)
				}
				ids = append(ids, q.Id)
			}

			// Check that answer lenghts are correct
			for i, q := range got.Questions {
				if len(tt.entryQuestions[i].Answer) != q.Length {
					t.Errorf("Unexpected answer length, got %v, want %v", q.Length, len(tt.entryQuestions[i].Answer))
				}
			}
		})
	}
}

func TestEntryQuiz_VerifyAnswers(t *testing.T) {

	ctxNoUser := context.Background()
	ctxSignedIn := ctxval.ContextWithUser(context.Background(), model.User{ID: "f00"})

	service1 := NewEntryQuizService([]config.EntryQuestion{
		{
			Question: "the best place to find a fox",
			Answer:   "panel",
		},
	})

	service3 := NewEntryQuizService([]config.EntryQuestion{
		{
			Question: "the best place to find a fox",
			Answer:   "woods",
		},
		{
			Question: "How many chickens can lay an egg?",
			Answer:   "Three",
		},
		{
			Question: "The US's first president",
			Answer:   "washington",
		},
	})

	tests := []struct {
		name       string
		eq         EntryQuiz
		ctx        context.Context
		answers    map[int]string
		wantPassed bool
		wantErr    bool
	}{
		{
			name:       "1 Quesion - no answers",
			ctx:        ctxSignedIn,
			eq:         service1,
			answers:    map[int]string{},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "1 Quesion - wrong",
			ctx:  ctxSignedIn,
			eq:   service1,
			answers: map[int]string{
				0: "ocean",
			},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "1 Quesion - too many answers",
			ctx:  ctxSignedIn,
			eq:   service1,
			answers: map[int]string{
				0: "ocean",
				1: "claimant",
			},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "1 Quesion - right answer",
			ctx:  ctxSignedIn,
			eq:   service1,
			answers: map[int]string{
				0: "PANeL",
			},
			wantPassed: true,
			wantErr:    false,
		},
		{
			name:       "3 Quesion - no answer",
			ctx:        ctxSignedIn,
			eq:         service3,
			answers:    map[int]string{},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "3 Quesion - not enough answers",
			ctx:  ctxSignedIn,
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				2: "washington",
			},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "3 Quesion - wrong 1",
			ctx:  ctxSignedIn,
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				1: "chair",
				2: "washington",
			},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "3 Quesion - too many answers",
			ctx:  ctxSignedIn,
			eq:   service3,
			answers: map[int]string{
				0: "ocean",
				1: "claimant",
				2: "legged",
				3: "RANGES",
			},
			wantPassed: false,
			wantErr:    false,
		},
		{
			name: "3 Quesion - right",
			ctx:  ctxSignedIn,
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				1: "THREE",
				2: "washingTON",
			},
			wantPassed: true,
			wantErr:    false,
		},
		{
			name: "3 Quesion - right, not signed in",
			ctx:  ctxNoUser,
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				1: "THREE",
				2: "washingTON",
			},
			wantPassed: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPassed, err := tt.eq.VerifyAnswers(tt.ctx, tt.answers)
			if (err != nil) != tt.wantErr {
				t.Errorf("EntryQuiz.VerifyAnswers() unexpected error value")
			}
			if gotPassed != tt.wantPassed {
				t.Errorf("EntryQuiz.VerifyAnswers() = %v, want %v", gotPassed, tt.wantPassed)
			}
		})
	}
}
