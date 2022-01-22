package service

import "strings"

// QuizQuestion is a crossword style question presented to the user to verify them before
// gaining access to the application.
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
	entryQuestions []QuizQuestion
}

// NewEntryQuizService creates and initialzies an EntryQuizService.
func NewEntryQuizService(entryQuestions []QuizQuestion) EntryQuiz {
	quiz := EntryQuiz{
		entryQuestions: entryQuestions,
	}

	for i := range quiz.entryQuestions {
		q := &quiz.entryQuestions[i]
		q.Id = i
		q.Length = len(q.Answer)
	}

	return quiz
}

// VerifyAnswers accepts a map of question IDs and string responses, and checks them
// against the correct answers, returning true if all correct, and false otherwise.
func (eq EntryQuiz) VerifyAnswers(answers map[int]string) bool {
	for _, q := range eq.entryQuestions {
		if strings.ToLower(q.Answer) != strings.ToLower(answers[q.Id]) {
			return false
		}
	}
	return true
}
