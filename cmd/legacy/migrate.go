package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rs/xid"

	"github.com/willbicks/epigram/internal/model"
	"github.com/willbicks/epigram/internal/storage/sqlite"
)

type legacyQuote struct {
	Date  int
	Quote string
	Sayer string
	Title string
}

// isValidLegacyQuoteMap checks whether the provided map has elements, and if each one has a quote,
// and returns the validity as a bool.
func isValidLegacyQuoteMap(qs map[string]legacyQuote) bool {
	if len(qs) < 1 {
		return false
	}
	for _, v := range qs {
		if v.Quote == "" {
			return false
		}
	}
	return true
}

// findQuoteArray recursively iterates through a raw json message until it finds an array / map
// of legacyQuote elements. Once it is found, it is unmarshalled and returned. If none found,
// returns nil.
func findQuoteArray(jMsg json.RawMessage) map[string]legacyQuote {
	var root map[string]json.RawMessage

	if err := json.Unmarshal(jMsg, &root); err != nil {
		return nil
	}
	for _, v := range root {
		quotes := make(map[string]legacyQuote)
		json.Unmarshal(v, &quotes)

		if isValidLegacyQuoteMap(quotes) {
			return quotes
		}

		quotes = findQuoteArray(v)
		if isValidLegacyQuoteMap(quotes) {
			return quotes
		}
	}
	return nil
}

// migrateQuote takes a legacy quote, normalizes the timestamp precision, and migrates it
// to an epigram model.Quote. Returns an error if the quote cannot be stored by repository.
func migrateQuote(repo *sqlite.QuoteRepository, uid string, q legacyQuote) error {
	// Normalize unix timestamp precision. If number of digits is less than 12, it is likely
	// in seconds format instead of millisecond, and should be multiplied by 1000 to convert.
	if q.Date < 1000000000000 {
		q.Date = q.Date * 1000
	}

	new := model.Quote{
		ID:          xid.New().String(),
		SubmitterID: uid,

		Quotee:  q.Sayer,
		Context: q.Title,
		Quote:   q.Quote,

		Created: time.UnixMilli(int64(q.Date)),
	}
	return repo.Create(context.Background(), new)
}
