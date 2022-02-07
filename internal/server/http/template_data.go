package http

// TemplateData is a struct which contains global site variables, as well as a page member,
// used for page specific data.
type TemplateData struct {
	Title string
	Page  interface{}
}

// joinPage returns a new TemplateData struct with the default site variables,
// and the specified page data object.
func (td TemplateData) joinPage(pd interface{}) TemplateData {
	td.Page = pd
	return td
}
