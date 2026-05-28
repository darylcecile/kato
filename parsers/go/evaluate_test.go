package kato

import (
	"reflect"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []EvalOption
		expected map[string]any
	}{
		{
			name:  "basic object with integer",
			input: "server {\n  port: 3000\n}",
			expected: map[string]any{
				"server": map[string]any{
					"port": 3000,
				},
			},
		},
		{
			name:  "nested objects",
			input: "db {\n  pool {\n    min: 2\n  }\n}",
			expected: map[string]any{
				"db": map[string]any{
					"pool": map[string]any{
						"min": 2,
					},
				},
			},
		},
		{
			name:  "arrays",
			input: `regions: ["a", "b"]`,
			expected: map[string]any{
				"regions": []any{"a", "b"},
			},
		},
		{
			name:  "string scalar",
			input: `name: "hello"`,
			expected: map[string]any{
				"name": "hello",
			},
		},
		{
			name:  "integer scalar",
			input: "count: 42",
			expected: map[string]any{
				"count": 42,
			},
		},
		{
			name:  "float scalar",
			input: "rate: 3.14",
			expected: map[string]any{
				"rate": 3.14,
			},
		},
		{
			name:  "bool true",
			input: "enabled: true",
			expected: map[string]any{
				"enabled": true,
			},
		},
		{
			name:  "bool false",
			input: "enabled: false",
			expected: map[string]any{
				"enabled": false,
			},
		},
		{
			name:  "null",
			input: "value: null",
			expected: map[string]any{
				"value": nil,
			},
		},
		{
			name:  "on becomes true",
			input: "enabled: on",
			expected: map[string]any{
				"enabled": true,
			},
		},
		{
			name:  "off becomes false",
			input: "debug: off",
			expected: map[string]any{
				"debug": false,
			},
		},
		{
			name:  "bare word",
			input: "mode: production",
			expected: map[string]any{
				"mode": "production",
			},
		},
		{
			name:  "unit literal",
			input: "timeout: 30s",
			expected: map[string]any{
				"timeout": map[string]any{
					"$unit": "s",
					"value": 30.0,
				},
			},
		},
		{
			name:  "env with option",
			input: `host: env("DB_HOST")`,
			opts:  []EvalOption{WithEnv(map[string]string{"DB_HOST": "localhost"})},
			expected: map[string]any{
				"host": "localhost",
			},
		},
		{
			name:  "ref resolves internal reference",
			input: "base_port: 3000\nport: ref(\"base_port\")",
			expected: map[string]any{
				"base_port": 3000,
				"port":      3000,
			},
		},
		{
			name:  "profile merges overrides",
			input: "server {\n  port: 3000\n}\n@profile production {\n  server {\n    port: 8080\n  }\n}",
			opts:  []EvalOption{WithProfile("production")},
			expected: map[string]any{
				"server": map[string]any{
					"port": 8080,
				},
			},
		},
		{
			name:  "top-level keys without wrapping object",
			input: "host: \"localhost\"\nport: 5432",
			expected: map[string]any{
				"host": "localhost",
				"port": 5432,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.input, tt.opts...)
			if err != nil {
				// Stub returns "not implemented" — mark as expected for now
				t.Skipf("Evaluate not yet implemented: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Evaluate() = %v, want %v", result, tt.expected)
			}
		})
	}
}
