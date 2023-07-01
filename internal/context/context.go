package context

import (
	"context"
	"crypto/rand"
	"encoding/hex"
)

const (
	placeHolder = "-"
)

// traceIDKey type is used for context value. The key is kept unexported to prevent collisions with keys defined in other packages.
type traceIDKey struct{}

func NewContextWithTraceID() context.Context {
	ctx := context.Background()
	traceID := make([]byte, 16) // Generate a 128 bit (16 bytes) trace ID
	_, err := rand.Read(traceID)
	if err != nil {
		return ctx
	}

	ctx = context.WithValue(ctx, traceIDKey{}, hex.EncodeToString(traceID))
	return ctx
}

// getTraceID returns the trace ID from the context. If context is nil or trace ID does not exist, it returns "-".
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return placeHolder
	}

	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return placeHolder
	}

	return traceID
}
