package frontend

import (
	"github.com/willbicks/epigram/internal/server/http/paths"
)

// RootTD is a struct which contains global site variables, as well as a page member,
// used for page specific data.
type RootTD struct {
	Title       string
	Description string
	Paths       paths.Paths
	Page        Page
}

// joinPage returns a new RootTD with the provided page joined to it.
func (td RootTD) joinPage(p Page) RootTD {
	td.Page = p
	return td
}
