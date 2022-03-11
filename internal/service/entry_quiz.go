package service

import (
	"context"
	"github.com/willbicks/epigram/internal/config"
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

// EntryQuiz contains a series of QuizQuestions new users must correctly answer before
// gaining access to the application.
type EntryQuiz struct {
	Questions []QuizQuestion
}

// NewEntryQuizService creates and initializes an EntryQuizService.
func NewEntryQuizService(qs []config.EntryQuestion) EntryQuiz {
	quiz := EntryQuiz{
		Questions: make([]QuizQuestion, len(qs)),
	}

	for i, q := range qs {
		quiz.Questions[i] = QuizQuestion{
			Id:       i,
			Length:   len(q.Answer),
			Question: q.Question,
			Answer:   q.Answer,
		}
	}

	return quiz
}

// VerifyAnswers accepts a map of question IDs and string responses, and checks them
// against the correct answers, and returns whether or not they match the expectation.
func (eq EntryQuiz) VerifyAnswers(ctx context.Context, answers map[int]string) (passed bool, err error) {

	if err := verifySignedIn(ctx); err != nil {
		return false, err
	}

	var wrongAnswer bool
	for _, q := range eq.Questions {
		if !strings.EqualFold(q.Answer, answers[q.Id]) {
			wrongAnswer = true
		}
	}

	return !wrongAnswer, nil
}
