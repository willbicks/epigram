package paths

// Paths stores url paths to each page to prevent hardcoding paths in
// multiple places.
type Paths struct {
	Home   string
	Quotes string
	Quiz   string
	Login  string
}

// Default returns the default paths assignments to be used in the application
func Default() Paths {
	return Paths{
		Home:   "/",
		Quotes: "/quotes",
		Quiz:   "/quiz",
		Login:  "/login",
	}
}
