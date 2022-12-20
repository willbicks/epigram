package frontend

import (
	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
)

// Page represents a page that can be rendered by the template engine, and any data to be inserted (if any).
type Page interface {
	viewName() string
}

// HomePage presents the home page
type HomePage struct {
}

func (HomePage) viewName() string {
	return "home.gohtml"
}

// PrivacyPage presents the privacy policy page
type PrivacyPage struct {
}

func (PrivacyPage) viewName() string {
	return "privacy.gohtml"
}

// QuotesPage lists all quotes by year
type QuotesPage struct {
	// RenderAdmin is true if the page should render admin controls / info
	RenderAdmin bool

	Error  error
	Quote  model.Quote
	Quotes []model.Quote

	// Users is a map of user ID to user, and should only be populated if RenderAdmin is true
	Users map[string]model.User
}

func (QuotesPage) viewName() string {
	return "quotes.gohtml"
}

// QuoteEditPage presents a form to edit a quote
type QuoteEditPage struct {
	Error error
	Quote model.Quote
}

func (QuoteEditPage) viewName() string {
	return "quote_edit.gohtml"
}

// QuizPage presents a quiz (list of questions)
type QuizPage struct {
	Error        error
	NumQuestions int
	Questions    []service.QuizQuestion
}

func (QuizPage) viewName() string {
	return "quiz.gohtml"
}

// AdminMainPage lists the users
type AdminMainPage struct {
	Users []model.User
}

func (AdminMainPage) viewName() string {
	return "admin_main.gohtml"
}
