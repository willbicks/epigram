package service

import (
	"testing"
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
		entryQuestions []QuizQuestion
	}{
		{
			name:           "No questions",
			entryQuestions: []QuizQuestion{},
		},
		{
			name: "One question",
			entryQuestions: []QuizQuestion{
				{
					Question: "How many chickens can lay an egg?",
					Answer:   "Three",
				},
			},
		},
		{
			name: "Three questions",
			entryQuestions: []QuizQuestion{
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

	service1 := NewEntryQuizService([]QuizQuestion{
		{
			Question: "the best place to find a fox",
			Answer:   "panel",
		},
	})

	service3 := NewEntryQuizService([]QuizQuestion{
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
		answers    map[int]string
		wantPassed bool
	}{
		{
			name:       "1 Quesion - no answers",
			eq:         service1,
			answers:    map[int]string{},
			wantPassed: false,
		},
		{
			name: "1 Quesion - wrong",
			eq:   service1,
			answers: map[int]string{
				0: "ocean",
			},
			wantPassed: false,
		},
		{
			name: "1 Quesion - too many answers",
			eq:   service1,
			answers: map[int]string{
				0: "ocean",
				1: "claimant",
			},
			wantPassed: false,
		},
		{
			name: "1 Quesion - right answer",
			eq:   service1,
			answers: map[int]string{
				0: "PANeL",
			},
			wantPassed: true,
		},
		{
			name:       "3 Quesion - no answer",
			eq:         service3,
			answers:    map[int]string{},
			wantPassed: false,
		},
		{
			name: "3 Quesion - not enough answers",
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				2: "washington",
			},
			wantPassed: false,
		},
		{
			name: "3 Quesion - wrong 1",
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				1: "chair",
				2: "washington",
			},
			wantPassed: false,
		},
		{
			name: "3 Quesion - too many answers",
			eq:   service3,
			answers: map[int]string{
				0: "ocean",
				1: "claimant",
				2: "legged",
				3: "RANGES",
			},
			wantPassed: false,
		},
		{
			name: "3 Quesion - right",
			eq:   service3,
			answers: map[int]string{
				0: "woods",
				1: "THREE",
				2: "washingTON",
			},
			wantPassed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPassed := tt.eq.VerifyAnswers(tt.answers); gotPassed != tt.wantPassed {
				t.Errorf("EntryQuiz.VerifyAnswers() = %v, want %v", gotPassed, tt.wantPassed)
			}
		})
	}
}
