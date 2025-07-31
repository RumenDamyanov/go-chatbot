package config

import "errors"

// Configuration validation errors.
var (
	ErrInvalidModel       = errors.New("invalid model specified")
	ErrInvalidTimeout     = errors.New("timeout must be greater than 0")
	ErrInvalidMaxTokens   = errors.New("max_tokens must be greater than 0")
	ErrInvalidTemperature = errors.New("temperature must be between 0 and 2")
	ErrMissingAPIKey      = errors.New("API key is required for this model")
	ErrMissingEndpoint    = errors.New("endpoint is required for this model")
	ErrUnsupportedModel   = errors.New("unsupported model")
)
