package paths

// Paths stores url paths to each page to prevent hardcoding paths in
// multiple places.
type Paths struct {
	Home      string
	Quotes    string
	QuoteEdit string
	Quiz      string
	Login     string
	Privacy   string
	Admin     string
}

// Default returns the default paths assignments to be used in the application
func Default() Paths {
	return Paths{
		Home:      "/",
		Quotes:    "/quotes",
		QuoteEdit: "/quote/edit",
		Quiz:      "/quiz",
		Login:     "/login",
		Privacy:   "/privacy",
		Admin:     "/admin",
	}
}
