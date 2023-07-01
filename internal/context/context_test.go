package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetTraceID_HappyCase(t *testing.T) {
	ctx := NewContextWithTraceID()
	assert.NotEqual(t, placeHolder, GetTraceID(ctx))
}

func TestGetTraceIDFromNil_HappyCase(t *testing.T) {
	var ctx context.Context
	assert.Equal(t, placeHolder, GetTraceID(ctx))
}
