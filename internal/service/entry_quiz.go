package service

import (
	"github.com/willbicks/charisms/internal/model"
	"strings"
)

// QuizQuestion is a crossword style question presented to the user to verify them before
// gaining access to the http.
type QuizQuestion struct {
	Id int
	// Length contains the number of characters in the Answer
	Length   int
	Question string
	Answer   string
}

// EntryQuiz contains a series of entryQuestions new users must correctly answer before
// gaining access to the applicaiton.
type EntryQuiz struct {
	EntryQuestions []QuizQuestion
}

// NewEntryQuizService creates and initialzies an EntryQuizService.
func NewEntryQuizService(entryQuestions []QuizQuestion) EntryQuiz {
	quiz := EntryQuiz{
		EntryQuestions: entryQuestions,
	}

	for i := range quiz.EntryQuestions {
		q := &quiz.EntryQuestions[i]
		q.Id = i
		q.Length = len(q.Answer)
	}

	return quiz
}

// VerifyAnswers accepts a map of question IDs and string responses, and checks them
// against the correct answers It increments the user's attempt counter, and records
// whether or not the user passed in u.QuizPassed.
func (eq EntryQuiz) VerifyAnswers(answers map[int]string, u *model.User) {
	u.QuizAttempts++
	for _, q := range eq.EntryQuestions {
		if strings.ToLower(q.Answer) != strings.ToLower(answers[q.Id]) {
			u.QuizPassed = false
		}
	}
	u.QuizPassed = true
}
