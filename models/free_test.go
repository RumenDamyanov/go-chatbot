package models

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFreeModel(t *testing.T) {
	model := NewFreeModel()

	assert.NotNil(t, model)
	assert.Equal(t, "free-model", model.Name())
	assert.Equal(t, "local", model.Provider())
}

func TestFreeModelAsk(t *testing.T) {
	model := NewFreeModel()
	ctx := context.Background()

	tests := []struct {
		name           string
		message        string
		expectedSubstr string
	}{
		{
			name:           "hello message",
			message:        "hello",
			expectedSubstr: "Hello",
		},
		{
			name:           "how are you message",
			message:        "how are you",
			expectedSubstr: "doing well",
		},
		{
			name:           "thank you message",
			message:        "thank you",
			expectedSubstr: "welcome",
		},
		{
			name:           "goodbye message",
			message:        "goodbye",
			expectedSubstr: "Goodbye",
		},
		{
			name:           "help message",
			message:        "help me",
			expectedSubstr: "help",
		},
		{
			name:           "name question",
			message:        "what's your name",
			expectedSubstr: "Bot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := model.Ask(ctx, tt.message, nil)

			assert.NoError(t, err)
			assert.NotEmpty(t, response)
			assert.Contains(t, response, tt.expectedSubstr)
		})
	}
}

func TestFreeModelHealth(t *testing.T) {
	model := NewFreeModel()
	ctx := context.Background()

	err := model.Health(ctx)
	assert.NoError(t, err)
}

func TestFreeModelCancellation(t *testing.T) {
	model := NewFreeModel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	response, err := model.Ask(ctx, "hello", nil)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Empty(t, response)
}

func TestFreeModelTimeout(t *testing.T) {
	model := NewFreeModel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// The free model has a 100ms delay, so this should timeout
	response, err := model.Ask(ctx, "hello", nil)

	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Empty(t, response)
}
