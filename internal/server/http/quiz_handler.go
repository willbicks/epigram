package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/willbicks/epigram/internal/ctxval"
	"github.com/willbicks/epigram/internal/service"
)

// quizTD represents the template data (TD) needed to render the quiz page
type quizTD struct {
	Error        error
	NumQuestions int
	Questions    []service.QuizQuestion
}

// quizHandler handles requests to the quizPage, either GET requests to render the page,
// or POST requests to submit attempts.
func (s *QuoteServer) quizHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := s.renderPage(w, "quiz.gohtml", quizTD{
			NumQuestions: len(s.QuizService.Questions),
			Questions:    s.QuizService.Questions,
		})
		if err != nil {
			s.serverError(w, r, err)
			return
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			s.clientError(w, r, err, http.StatusBadRequest)
			return
		}

		answers := make(map[int]string)
		for id, value := range r.PostForm {
			id, err := strconv.Atoi(id)
			if err != nil {
				s.clientError(w, r, errors.New("invalid form value key"), http.StatusBadRequest)
				return
			}
			answers[id] = value[0]
		}

		u := ctxval.UserFromContext(r.Context())
		passed, err := s.QuizService.VerifyAnswers(r.Context(), answers)
		if err != nil {
			s.serverError(w, r, err)
		}

		failReason, err := s.UserService.RecordQuizAttempt(r.Context(), &u, passed)
		if err != nil {
			s.serverError(w, r, err)
		}

		if u.QuizPassed {
			http.Redirect(w, r, s.Config.routes.Quotes, http.StatusFound)
			return
		} else {
			err := s.renderPage(w, "quiz.gohtml", quizTD{
				NumQuestions: len(s.QuizService.Questions),
				Questions:    s.QuizService.Questions,
				Error:        errors.New(failReason),
			})
			if err != nil {
				s.serverError(w, r, err)
				return
			}
		}

	default:
		s.methodNotAllowedError(w, r)
		return
	}
}
