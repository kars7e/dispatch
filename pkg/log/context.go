package log

import (
	"context"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

type (
	logKey struct{}
)

var (
	stdLog = logrus.NewEntry(logrus.StandardLogger())
)

// WithLogger returns a new context with the provided logger.
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, logKey{}, logger)
}

// GetLogger retrieves the current logger from the context. If no logger is
// available, the default logger is returned.
func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(logKey{})

	if logger == nil {
		return stdLog
	}

	return logger.(*logrus.Entry)
}

// WithRequestID returns a logger with requestID and a context with this logger injected
func WithRequestID(ctx context.Context) (*logrus.Entry, context.Context) {
	val := ctx.Value(logKey{})

	if val == nil {
		val = stdLog.WithField("requestId", uuid.NewV4().String())
	}
	logger := val.(*logrus.Entry)

	return logger, WithLogger(ctx, logger)
}
