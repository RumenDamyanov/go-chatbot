package models

import (
	"fmt"
	"regexp"
)

// safeModelName matches ids that can be placed in a URL path or a JSON field
// without changing the request target.
var safeModelName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// normalizeModel uses fallback when configured is empty, and rejects names
// that contain separators or control characters.
func normalizeModel(configured, fallback string) (string, error) {
	model := configured
	if model == "" {
		model = fallback
	}
	if !safeModelName.MatchString(model) {
		return "", fmt.Errorf("the model name contains characters that are not allowed")
	}
	return model, nil
}
