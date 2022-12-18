package frontend

import (
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/server/http/paths"
	"github.com/willbicks/epigram/internal/service"
)

// RootTD is a struct which contains global site variables, as well as a page member,
// used for page specific data.
type RootTD struct {
	Title       string
	Description string
	Paths       paths.Paths
	Page        interface{}
}

// joinPage returns a new TemplateData struct with the default site variables,
// and the specified page data object.
func (td RootTD) joinPage(pd interface{}) RootTD {
	td.Page = pd
	return td
}

// QuotesTD represents the template data (TD) needed to render the quotes page
type QuotesTD struct {
	Error  error
	Quote  model.Quote
	Quotes []model.Quote
}

// QuizTD represents the template data (TD) needed to render the quiz page
type QuizTD struct {
	Error        error
	NumQuestions int
	Questions    []service.QuizQuestion
}

// AdminMainTD is the data needed to render the main admin template.
type AdminMainTD struct {
	Users []model.User
}
