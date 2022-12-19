package frontend

import (
	"bytes"
	"strings"
	"testing"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
)

func Test_NewTemplateEngine(t *testing.T) {
	_, err := NewTemplateEngine(RootTD{})
	if err != nil {
		t.Error("NewTemplateEngine() returned error:", err)
	}
}
func Test_TemplateEngine_RenderPage(t *testing.T) {
	tests := []Page{
		HomePage{},
		PrivacyPage{},
		QuotesPage{
			Quotes: []model.Quote{
				{
					Quotee:  "Test Quotee",
					Quote:   "Test Quote",
					Context: "Test Context",
				},
			},
		},
		QuotesPage{
			RenderAdmin: true,
			Quotes: []model.Quote{
				{
					Quotee:      "Test Quotee",
					Quote:       "Test Quote",
					Context:     "Test Context",
					SubmitterID: "x123",
				},
			},
			Users: map[string]model.User{
				"x123": {
					ID:   "x123",
					Name: "Test User",
				},
			},
		},
		QuizPage{
			Questions: []service.QuizQuestion{
				{
					ID:       1,
					Question: "Test Question",
					Answer:   "Answer",
					Length:   6,
				},
			},
		},
		AdminMainPage{
			Users: []model.User{
				{
					ID:    "x123",
					Name:  "Test User",
					Email: "test@example.com",
				},
			},
		},
	}

	te, err := NewTemplateEngine(RootTD{})
	if err != nil {
		t.Error("NewTemplateEngine() returned error:", err)
	}

	for _, p := range tests {
		t.Run(p.viewName(), func(t *testing.T) {
			var buf bytes.Buffer
			err := te.RenderPage(&buf, p)
			if err != nil {
				t.Errorf("RenderPage() for %s returned error: %s", p.viewName(), err)
			}
			rendered := strings.TrimSpace(buf.String())
			if rendered[len(rendered)-7:] != "</html>" {
				t.Errorf("RenderPage() for %s does not appear to render complete HTML", p.viewName())
			}
		})
	}
}
