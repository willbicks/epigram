package frontend

import (
	"bytes"
	"strings"
	"testing"
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
		QuotesPage{},
		QuizPage{},
		AdminMainPage{},
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
