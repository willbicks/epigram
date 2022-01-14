package application

// TemplateData is a struct which contains global site variables, as well as a page member,
// used for page specific data.
type TemplateData struct {
	Title string
	page  interface{}
}

// joinPage returns a new TemplateData struct with the default site variables,
// and the specified page data object.
func (td TemplateData) joinPage(pd interface{}) TemplateData {
	td.page = pd
	return td
}
