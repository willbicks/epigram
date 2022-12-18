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
	tests := []struct {
		templateName string
		templateData interface{}
	}{
		{
			"home.gohtml",
			nil,
		},
		{
			"privacy.gohtml",
			nil,
		},
		{
			"admin_main.gohtml",
			AdminMainTD{},
		},
		{
			"quotes.gohtml",
			QuotesTD{},
		},
		{
			"quiz.gohtml",
			QuizTD{},
		},
	}

	te, err := NewTemplateEngine(RootTD{})
	if err != nil {
		t.Error("NewTemplateEngine() returned error:", err)
	}

	for _, tt := range tests {
		t.Run(tt.templateName, func(t *testing.T) {
			var buf bytes.Buffer
			err := te.RenderPage(&buf, tt.templateName, tt.templateData)
			if err != nil {
				t.Errorf("RenderPage() for %s returned error: %s", tt.templateName, err)
			}
			rendered := strings.TrimSpace(buf.String())
			if rendered[len(rendered)-7:] != "</html>" {
				t.Errorf("RenderPage() for %s does not appear to render complete HTML", tt.templateName)
			}
		})
	}
}
