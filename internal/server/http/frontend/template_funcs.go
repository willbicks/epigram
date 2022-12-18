package frontend

import (
	"fmt"
	"html/template"
	"net/url"
	"sort"
	"strings"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/service"
)

var templateFuncs template.FuncMap = template.FuncMap{
	// getIssues takes an error and returns a slice of issues. Derived from issues slice if ServiceError,
	// otherwise a slice with a single entry of the error.Error()
	"getIssues": func(err error) []string {
		if err == nil {
			return []string{}
		}

		serr, ok := err.(service.Error)
		if !ok {
			return []string{err.Error()}
		}

		return serr.Issues
	},
	// quotesByYear takes a slice of quotes, and returns them as a map where the key is the year.
	"quotesByYear": func(quotes []model.Quote) map[int][]model.Quote {
		// sort quotes from newest to oldest
		sort.Slice(quotes, func(i, j int) bool {
			return quotes[i].Created.After(quotes[j].Created)
		})

		byYear := make(map[int][]model.Quote)

		for _, q := range quotes {
			byYear[q.Created.Year()] = append(byYear[q.Created.Year()], q)
		}

		return byYear
	},
	// orderedYearKeys takes a map of quotes, and returns a slice of years (keys of map) in
	// reverse chronological order
	"orderedYearKeys": func(quotes map[int][]model.Quote) []int {
		years := make([]int, len(quotes))
		i := 0
		for y := range quotes {
			years[i] = y
			i++
		}
		sort.Slice(years, func(i, j int) bool {
			return years[i] > years[j]
		})
		return years
	},
	// sizeImage accepts a url of an image, and attempts to resize it by modifying the urlparams of the url,
	// depending on the image hosting service. Currently supports googleusercontent. Returns the url of the
	// sizedImage, or returns a url to a not found image placeholder if not a valid URL.
	"sizeImage": func(imgURL string, size int) string {
		if imgURL == "" {
			return "/static/img/notfound.png"
		}

		u, err := url.Parse(imgURL)
		if err != nil {
			return "/static/img/notfound.png"
		}

		switch u.Hostname() {
		case "lh3.googleusercontent.com":
			imgURL := imgURL[:strings.LastIndex(imgURL, "=")]
			return fmt.Sprintf("%s=s%d-c", imgURL, size)
		}

		return imgURL
	},
}
