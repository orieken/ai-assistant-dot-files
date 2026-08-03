package logging

import (
	"io"
	"log/slog"
)

// Logger is a thin structured-logging wrapper.
// Tools receive *Logger by constructor injection; the handler wires it from os.Stderr.
type Logger struct {
	inner *slog.Logger
}

// NewLogger constructs a Logger that writes JSON-encoded structured records to w.
func NewLogger(w io.Writer) *Logger {
	return &Logger{inner: slog.New(slog.NewJSONHandler(w, nil))}
}

func (l *Logger) Info(msg string, args ...any)  { l.inner.Info(msg, args...) }
func (l *Logger) Warn(msg string, args ...any)  { l.inner.Warn(msg, args...) }
func (l *Logger) Error(msg string, args ...any) { l.inner.Error(msg, args...) }
