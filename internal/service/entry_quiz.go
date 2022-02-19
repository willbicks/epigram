package service

import (
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
// gaining access to the applicaiton.
type EntryQuiz struct {
	Questions []QuizQuestion
}

// NewEntryQuizService creates and initialzies an EntryQuizService.
func NewEntryQuizService(entryQuestions []QuizQuestion) EntryQuiz {
	quiz := EntryQuiz{
		Questions: entryQuestions,
	}

	for i := range quiz.Questions {
		q := &quiz.Questions[i]
		q.Id = i
		q.Length = len(q.Answer)
	}

	return quiz
}

// VerifyAnswers accepts a map of question IDs and string responses, and checks them
// against the correct answers, and returns whether or not they match the expectation.
func (eq EntryQuiz) VerifyAnswers(answers map[int]string) (passed bool) {
	var wrongAnswer bool
	for _, q := range eq.Questions {
		if !strings.EqualFold(q.Answer, answers[q.Id]) {
			wrongAnswer = true
		}
	}

	return !wrongAnswer
}
