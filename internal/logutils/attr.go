package logutils

import "log/slog"

// Error returns an Attr for an error with the key "err".
func Error(err error) slog.Attr {
	return slog.String("err", err.Error())
}
