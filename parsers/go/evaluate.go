package kato

import "fmt"

// EvalOption configures the Evaluate function.
type EvalOption func(*evalConfig)

type evalConfig struct {
	env     map[string]string
	profile string
}

// WithEnv provides environment variables for resolving env() calls.
func WithEnv(env map[string]string) EvalOption {
	return func(c *evalConfig) {
		c.env = env
	}
}

// WithProfile selects a named profile whose overrides are merged into the result.
func WithProfile(name string) EvalOption {
	return func(c *evalConfig) {
		c.profile = name
	}
}

// Evaluate parses Kato input and returns a plain map representation.
//
// Conversion rules:
//   - Objects become map[string]any
//   - Arrays become []any
//   - on → true, off → false
//   - Unit literals (e.g. 30s) → map[string]any{"$unit": "<suffix>", "value": <number>}
//   - env() resolves using WithEnv option
//   - ref() resolves internal references
//   - Profiles are merged when WithProfile is used
func Evaluate(input string, opts ...EvalOption) (map[string]any, error) {
	_ = &evalConfig{}
	for _, o := range opts {
		o(&evalConfig{})
	}
	return nil, fmt.Errorf("not implemented")
}
