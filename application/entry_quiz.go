package application

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/willbicks/charisms/service"
)

// quizTD represents the template data (TD) needed to render the quiz page
type quizTD struct {
	Issues       []string
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		answers := make(map[int]string)
		for id, value := range r.PostForm {
			id, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, "Invalid form value key", http.StatusBadRequest)
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
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}
}
