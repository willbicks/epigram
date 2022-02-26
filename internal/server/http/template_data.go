package http

// RootTD is a struct which contains global site variables, as well as a page member,
// used for page specific data.
type RootTD struct {
	Title       string
	Description string
	Routes      routeStruct
	Page        interface{}
}

// joinPage returns a new TemplateData struct with the default site variables,
// and the specified page data object.
func (td RootTD) joinPage(pd interface{}) RootTD {
	td.Page = pd
	return td
}
