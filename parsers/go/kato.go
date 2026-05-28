// Package kato implements a parser for the Kato configuration language.
package kato

import "fmt"

// Token types produced by the lexer.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNewline
	TokenComment
	TokenDocComment
	TokenDirective
	TokenIdentifier
	TokenString
	TokenInteger
	TokenFloat
	TokenBoolTrue
	TokenBoolFalse
	TokenNull
	TokenOn
	TokenOff
	TokenBareWord
	TokenLBrace
	TokenRBrace
	TokenLBracket
	TokenRBracket
	TokenColon
	TokenComma
	TokenDot
	TokenAt
	TokenPlusEq
	TokenMinusEq
	TokenLParen
	TokenRParen
	TokenSlash
	TokenUnitSuffix
)

// Token represents a single lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer tokenizes Kato source input.
type Lexer struct {
	input   string
	pos     int
	line    int
	col     int
	tokens  []Token
}

// NewLexer creates a new Lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{input: input, line: 1, col: 1}
}

// Tokenize performs lexical analysis and returns all tokens.
// TODO: implement full tokenization.
func (l *Lexer) Tokenize() ([]Token, error) {
	return nil, fmt.Errorf("lexer not yet implemented")
}

// --- AST Node Types ---

// NodeType identifies the kind of AST node.
type NodeType int

const (
	NodeDocument NodeType = iota
	NodeObject
	NodeArray
	NodeScalar
	NodeDirective
	NodeComment
	NodeDocComment
	NodeUnitLiteral
	NodeFunctionCall
	NodeKeyValue
	NodeArrayOp
)

// Node is the interface all AST nodes implement.
type Node interface {
	Type() NodeType
	Pos() Position
}

// Position in source.
type Position struct {
	Line   int
	Column int
}

// Document is the root AST node.
type Document struct {
	Position   Position
	Statements []Node
	Directives []*DirectiveNode
	Comments   []*CommentNode
}

func (d *Document) Type() NodeType { return NodeDocument }
func (d *Document) Pos() Position  { return d.Position }

// ObjectNode represents a Kato object { ... }.
type ObjectNode struct {
	Position Position
	Members  []Node
}

func (n *ObjectNode) Type() NodeType { return NodeObject }
func (n *ObjectNode) Pos() Position  { return n.Position }

// ArrayNode represents a Kato array [ ... ].
type ArrayNode struct {
	Position Position
	Elements []Node
}

func (n *ArrayNode) Type() NodeType { return NodeArray }
func (n *ArrayNode) Pos() Position  { return n.Position }

// ScalarNode represents a scalar value.
type ScalarNode struct {
	Position Position
	Value    interface{}
	Raw      string
	Kind     ScalarKind
}

// ScalarKind identifies the scalar subtype.
type ScalarKind int

const (
	ScalarString ScalarKind = iota
	ScalarInt
	ScalarFloat
	ScalarBool
	ScalarNull
	ScalarOnOff
	ScalarBareWord
)

func (n *ScalarNode) Type() NodeType { return NodeScalar }
func (n *ScalarNode) Pos() Position  { return n.Position }

// DirectiveNode represents an @directive.
type DirectiveNode struct {
	Position Position
	Name     string
	Args     []Node
	Body     Node // optional object body (e.g. @profile production { ... })
}

func (n *DirectiveNode) Type() NodeType { return NodeDirective }
func (n *DirectiveNode) Pos() Position  { return n.Position }

// CommentNode represents a # comment.
type CommentNode struct {
	Position Position
	Text     string
}

func (n *CommentNode) Type() NodeType { return NodeComment }
func (n *CommentNode) Pos() Position  { return n.Position }

// DocCommentNode represents a /// doc comment.
type DocCommentNode struct {
	Position Position
	Text     string
}

func (n *DocCommentNode) Type() NodeType { return NodeDocComment }
func (n *DocCommentNode) Pos() Position  { return n.Position }

// UnitLiteralNode represents a value with a unit suffix (e.g. 30s, 512MiB).
type UnitLiteralNode struct {
	Position Position
	Value    float64
	Unit     string
	Raw      string
}

func (n *UnitLiteralNode) Type() NodeType { return NodeUnitLiteral }
func (n *UnitLiteralNode) Pos() Position  { return n.Position }

// FunctionCallNode represents env() or ref() calls.
type FunctionCallNode struct {
	Position Position
	Name     string
	Args     []Node
}

func (n *FunctionCallNode) Type() NodeType { return NodeFunctionCall }
func (n *FunctionCallNode) Pos() Position  { return n.Position }

// KeyValueNode represents key: value.
type KeyValueNode struct {
	Position Position
	Key      string
	Value    Node
}

func (n *KeyValueNode) Type() NodeType { return NodeKeyValue }
func (n *KeyValueNode) Pos() Position  { return n.Position }

// ArrayOpNode represents key += [...] or key -= [...].
type ArrayOpNode struct {
	Position Position
	Key      string
	Op       string // "+=" or "-="
	Value    *ArrayNode
}

func (n *ArrayOpNode) Type() NodeType { return NodeArrayOp }
func (n *ArrayOpNode) Pos() Position  { return n.Position }

// NamedBlockNode represents key { ... }.
type NamedBlockNode struct {
	Position Position
	Name     string
	Body     *ObjectNode
}

func (n *NamedBlockNode) Type() NodeType { return NodeObject }
func (n *NamedBlockNode) Pos() Position  { return n.Position }

// NamedArrayNode represents key [ ... ].
type NamedArrayNode struct {
	Position Position
	Name     string
	Body     *ArrayNode
}

func (n *NamedArrayNode) Type() NodeType { return NodeArray }
func (n *NamedArrayNode) Pos() Position  { return n.Position }

// Parse parses Kato source and returns a Document AST.
// This is currently a stub — returns an error until the parser is implemented.
func Parse(input string) (*Document, error) {
	return nil, fmt.Errorf("parser not yet implemented")
}
