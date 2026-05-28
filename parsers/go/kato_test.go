package kato

import (
	"strings"
	"testing"
)

// TestParse_BasicObject tests parsing simple object blocks.
func TestParse_BasicObject(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name:  "simple object with integer value",
			input: "server {\n  port: 3000\n}",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(doc.Statements) != 1 {
					t.Fatalf("expected 1 statement, got %d", len(doc.Statements))
				}
				nb, ok := doc.Statements[0].(*NamedBlockNode)
				if !ok {
					t.Fatalf("expected NamedBlockNode, got %T", doc.Statements[0])
				}
				if nb.Name != "server" {
					t.Fatalf("expected name 'server', got %q", nb.Name)
				}
				if len(nb.Body.Members) != 1 {
					t.Fatalf("expected 1 member, got %d", len(nb.Body.Members))
				}
				kv, ok := nb.Body.Members[0].(*KeyValueNode)
				if !ok {
					t.Fatalf("expected KeyValueNode, got %T", nb.Body.Members[0])
				}
				if kv.Key != "port" {
					t.Fatalf("expected key 'port', got %q", kv.Key)
				}
				scalar, ok := kv.Value.(*ScalarNode)
				if !ok {
					t.Fatalf("expected ScalarNode, got %T", kv.Value)
				}
				if scalar.Kind != ScalarInt || scalar.Value != int64(3000) {
					t.Fatalf("expected int 3000, got %v (%v)", scalar.Value, scalar.Kind)
				}
			},
		},
		{
			name:  "object with string value",
			input: `server { host: "localhost" }`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				nb := doc.Statements[0].(*NamedBlockNode)
				kv := nb.Body.Members[0].(*KeyValueNode)
				scalar := kv.Value.(*ScalarNode)
				if scalar.Kind != ScalarString || scalar.Value != "localhost" {
					t.Fatalf("expected string 'localhost', got %v", scalar.Value)
				}
			},
		},
		{
			name:  "object with multiple fields",
			input: "server {\n  port: 3000\n  host: \"localhost\"\n  enabled: true\n}",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				nb := doc.Statements[0].(*NamedBlockNode)
				if len(nb.Body.Members) != 3 {
					t.Fatalf("expected 3 members, got %d", len(nb.Body.Members))
				}
			},
		},
		{
			name:  "empty object",
			input: "server {}",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				nb := doc.Statements[0].(*NamedBlockNode)
				if len(nb.Body.Members) != 0 {
					t.Fatalf("expected 0 members, got %d", len(nb.Body.Members))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_NestedObjects tests parsing nested object structures.
func TestParse_NestedObjects(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name:  "single nested object",
			input: "database {\n  pool {\n    min: 2\n    max: 10\n  }\n}",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				db := doc.Statements[0].(*NamedBlockNode)
				if db.Name != "database" {
					t.Fatalf("expected 'database', got %q", db.Name)
				}
				pool := db.Body.Members[0].(*NamedBlockNode)
				if pool.Name != "pool" {
					t.Fatalf("expected 'pool', got %q", pool.Name)
				}
				if len(pool.Body.Members) != 2 {
					t.Fatalf("expected 2 members in pool, got %d", len(pool.Body.Members))
				}
			},
		},
		{
			name:  "deeply nested",
			input: "a {\n  b {\n    c {\n      val: 1\n    }\n  }\n}",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				a := doc.Statements[0].(*NamedBlockNode)
				b := a.Body.Members[0].(*NamedBlockNode)
				c := b.Body.Members[0].(*NamedBlockNode)
				kv := c.Body.Members[0].(*KeyValueNode)
				if kv.Key != "val" {
					t.Fatalf("expected 'val', got %q", kv.Key)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_InlineArrays tests inline array syntax.
func TestParse_InlineArrays(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name:  "string array",
			input: `regions: ["eu-west-1", "us-east-1"]`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				kv := doc.Statements[0].(*KeyValueNode)
				arr := kv.Value.(*ArrayNode)
				if len(arr.Elements) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
				}
			},
		},
		{
			name:  "integer array",
			input: `ports: [80, 443, 8080]`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				kv := doc.Statements[0].(*KeyValueNode)
				arr := kv.Value.(*ArrayNode)
				if len(arr.Elements) != 3 {
					t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
				}
			},
		},
		{
			name:  "trailing comma",
			input: `items: ["a", "b", "c",]`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				kv := doc.Statements[0].(*KeyValueNode)
				arr := kv.Value.(*ArrayNode)
				if len(arr.Elements) != 3 {
					t.Fatalf("expected 3 elements with trailing comma, got %d", len(arr.Elements))
				}
			},
		},
		{
			name:  "empty array",
			input: `items: []`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				kv := doc.Statements[0].(*KeyValueNode)
				arr := kv.Value.(*ArrayNode)
				if len(arr.Elements) != 0 {
					t.Fatalf("expected 0 elements, got %d", len(arr.Elements))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_BlockArrays tests block (newline-separated) arrays.
func TestParse_BlockArrays(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name:  "named block array",
			input: "plugins [\n  \"auth\"\n  \"payments\"\n]",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				na := doc.Statements[0].(*NamedArrayNode)
				if na.Name != "plugins" {
					t.Fatalf("expected 'plugins', got %q", na.Name)
				}
				if len(na.Body.Elements) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(na.Body.Elements))
				}
			},
		},
		{
			name:  "newline separated with trailing newline",
			input: "items [\n  \"a\"\n  \"b\"\n  \"c\"\n]",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				na := doc.Statements[0].(*NamedArrayNode)
				if len(na.Body.Elements) != 3 {
					t.Fatalf("expected 3 elements, got %d", len(na.Body.Elements))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_Scalars tests all scalar types.
func TestParse_Scalars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantKind ScalarKind
		wantVal  interface{}
	}{
		{"string", `name: "Daryl"`, ScalarString, "Daryl"},
		{"integer", `age: 25`, ScalarInt, int64(25)},
		{"negative integer", `offset: -10`, ScalarInt, int64(-10)},
		{"float", `ratio: 3.14`, ScalarFloat, 3.14},
		{"negative float", `temp: -0.5`, ScalarFloat, -0.5},
		{"bool true", `enabled: true`, ScalarBool, true},
		{"bool false", `enabled: false`, ScalarBool, false},
		{"null", `value: null`, ScalarNull, nil},
		{"on", `feature: on`, ScalarOnOff, true},
		{"off", `feature: off`, ScalarOnOff, false},
		{"bare word", `env: production`, ScalarBareWord, "production"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			kv := doc.Statements[0].(*KeyValueNode)
			scalar := kv.Value.(*ScalarNode)
			if scalar.Kind != tt.wantKind {
				t.Fatalf("expected kind %v, got %v", tt.wantKind, scalar.Kind)
			}
			if scalar.Value != tt.wantVal {
				t.Fatalf("expected value %v, got %v", tt.wantVal, scalar.Value)
			}
		})
	}
}

// TestParse_UnitLiterals tests unit-suffixed values.
func TestParse_UnitLiterals(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantVal  float64
		wantUnit string
	}{
		{"seconds", `timeout: 30s`, 30, "s"},
		{"mebibytes", `memory: 512MiB`, 512, "MiB"},
		{"rate", `rateLimit: 100/min`, 100, "/min"},
		{"milliseconds", `delay: 250ms`, 250, "ms"},
		{"gigabytes", `disk: 50GB`, 50, "GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			kv := doc.Statements[0].(*KeyValueNode)
			unit := kv.Value.(*UnitLiteralNode)
			if unit.Value != tt.wantVal {
				t.Fatalf("expected value %v, got %v", tt.wantVal, unit.Value)
			}
			if unit.Unit != tt.wantUnit {
				t.Fatalf("expected unit %q, got %q", tt.wantUnit, unit.Unit)
			}
		})
	}
}

// TestParse_Comments tests comment parsing.
func TestParse_Comments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name:  "line comment",
			input: "# this is a comment\nport: 3000",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				// Document should have the comment and the key-value
				if len(doc.Comments) < 1 {
					t.Fatal("expected at least 1 comment")
				}
			},
		},
		{
			name:  "doc comment",
			input: "/// This documents the port\nport: 3000",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				found := false
				for _, s := range doc.Statements {
					if dc, ok := s.(*DocCommentNode); ok {
						if strings.Contains(dc.Text, "documents the port") {
							found = true
						}
					}
				}
				if !found {
					t.Fatal("expected doc comment to be preserved")
				}
			},
		},
		{
			name:  "inline comment after value",
			input: "port: 3000  # the port number",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				// The key-value should still parse; comment is preserved separately
				if len(doc.Statements) < 1 {
					t.Fatal("expected at least 1 statement")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_Directives tests @directive parsing.
func TestParse_Directives(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name:  "@version",
			input: "@version 1",
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(doc.Directives) != 1 {
					t.Fatalf("expected 1 directive, got %d", len(doc.Directives))
				}
				if doc.Directives[0].Name != "version" {
					t.Fatalf("expected 'version', got %q", doc.Directives[0].Name)
				}
			},
		},
		{
			name:  "@schema",
			input: `@schema "https://example.com/schema.kato"`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if doc.Directives[0].Name != "schema" {
					t.Fatalf("expected 'schema', got %q", doc.Directives[0].Name)
				}
			},
		},
		{
			name:  "@include",
			input: `@include "./shared.kato"`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if doc.Directives[0].Name != "include" {
					t.Fatalf("expected 'include', got %q", doc.Directives[0].Name)
				}
			},
		},
		{
			name:  "@include with alias",
			input: `@include "./secrets.kato" as secrets`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if doc.Directives[0].Name != "include" {
					t.Fatalf("expected 'include', got %q", doc.Directives[0].Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_FunctionCalls tests env() and ref() parsing.
func TestParse_FunctionCalls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantFunc string
		wantArgs int
	}{
		{"env simple", `url: env("DATABASE_URL")`, "env", 1},
		{"env with fallback", `url: env("DB_URL", fallback: null)`, "env", 2},
		{"ref", `region: ref("shared.region")`, "ref", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			kv := doc.Statements[0].(*KeyValueNode)
			fc := kv.Value.(*FunctionCallNode)
			if fc.Name != tt.wantFunc {
				t.Fatalf("expected func %q, got %q", tt.wantFunc, fc.Name)
			}
			if len(fc.Args) != tt.wantArgs {
				t.Fatalf("expected %d args, got %d", tt.wantArgs, len(fc.Args))
			}
		})
	}
}

// TestParse_Profiles tests @profile directive with body.
func TestParse_Profiles(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, doc *Document, err error)
	}{
		{
			name: "production profile",
			input: `@profile production {
  server {
    workers: 8
  }
}`,
			check: func(t *testing.T, doc *Document, err error) {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(doc.Directives) != 1 {
					t.Fatalf("expected 1 directive, got %d", len(doc.Directives))
				}
				d := doc.Directives[0]
				if d.Name != "profile" {
					t.Fatalf("expected 'profile', got %q", d.Name)
				}
				if d.Body == nil {
					t.Fatal("expected profile to have a body")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			tt.check(t, doc, err)
		})
	}
}

// TestParse_ArrayOperations tests += and -= operators.
func TestParse_ArrayOperations(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantOp string
		wantN  int
	}{
		{"append", `plugins += ["analytics"]`, "+=", 1},
		{"remove", `plugins -= ["debug-toolbar"]`, "-=", 1},
		{"append multiple", `features += ["a", "b", "c"]`, "+=", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			op := doc.Statements[0].(*ArrayOpNode)
			if op.Op != tt.wantOp {
				t.Fatalf("expected op %q, got %q", tt.wantOp, op.Op)
			}
			if len(op.Value.Elements) != tt.wantN {
				t.Fatalf("expected %d elements, got %d", tt.wantN, len(op.Value.Elements))
			}
		})
	}
}

// TestParse_Errors tests that invalid syntax produces errors.
func TestParse_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unclosed brace", `server { port: 3000`},
		{"unclosed bracket", `items: ["a", "b"`},
		{"unclosed string", `name: "hello`},
		{"invalid directive", `@ 123`},
		{"colon without value", `key:`},
		{"double colon", `key:: value`},
		{"bare at sign", `@`},
		{"mismatched braces", `server { port: 3000 ]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

// TestParse_ComplexDocument tests a full document with mixed constructs.
func TestParse_ComplexDocument(t *testing.T) {
	input := `@version 1
@schema "https://example.com/app.kato"

# Application configuration
/// Main server settings
server {
  host: "0.0.0.0"
  port: 8080
  timeout: 30s
  workers: 4
}

database {
  url: env("DATABASE_URL")
  pool {
    min: 2
    max: 20
  }
}

regions: ["eu-west-1", "us-east-1"]

plugins [
  "auth"
  "payments"
  "search"
]

@profile production {
  server {
    workers: 16
  }
  plugins += ["monitoring"]
}
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have directives
	if len(doc.Directives) < 2 {
		t.Fatalf("expected at least 2 directives, got %d", len(doc.Directives))
	}

	// Should have multiple top-level statements
	if len(doc.Statements) < 4 {
		t.Fatalf("expected at least 4 statements, got %d", len(doc.Statements))
	}
}
