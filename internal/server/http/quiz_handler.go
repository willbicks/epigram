package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/willbicks/charisms/internal/service"
)

// quizTD represents the template data (TD) needed to render the quiz page
type quizTD struct {
	Error        error
	NumQuestions int
	Questions    []service.QuizQuestion
}

// quizHandler handles requests to the quizpage (/)
func (s *CharismsServer) quizHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := s.renderPage(w, "quiz.gohtml", quizTD{
			NumQuestions: len(s.QuizService.EntryQuestions),
			Questions:    s.QuizService.EntryQuestions,
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

		u := UserFromContext(r.Context())
		s.QuizService.VerifyAnswers(answers, &u)
		s.UserService.UpdateUser(r.Context(), u)

		if u.QuizPassed {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("answers correct"))
			return
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("at least one answer is incorrect"))
			return
		}

	default:
		s.methodNotAllowedError(w, r)
		return
	}
}