package model

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestQuote_Editable(t *testing.T) {
	joe := User{
		ID: "joe",
	}
	charlene := User{
		ID: "charlene",
	}
	admin := User{
		ID:    "admin",
		Admin: true,
	}

	newQuote := Quote{
		SubmitterID: joe.ID,
		Created:     time.Now(),
	}
	oldQuote := Quote{
		SubmitterID: joe.ID,
		Created:     time.Now().Add(-2 * time.Hour),
	}

	is := is.New(t)

	// Editable quote / user combinations
	is.True(newQuote.Editable(joe))
	is.True(newQuote.Editable(admin))
	is.True(oldQuote.Editable(admin))

	// Not editable quote / user combinations
	is.True(!newQuote.Editable(charlene))
	is.True(!oldQuote.Editable(joe))
	is.True(!oldQuote.Editable(charlene))
}
