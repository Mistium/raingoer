// Package declaration and imports
package ast

import "fmt"

// Node interface represents a generic AST node
type Node interface {
	String() string
}

// Program represents the root of the AST
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	result := "Program{\n"
	for _, stmt := range p.Statements {
		result += "  " + stmt.String() + "\n"
	}
	result += "}"
	return result
}

// Statement interface for all statement nodes
type Statement interface {
	Node
	statementNode()
}

// Expression interface for all expression nodes
type Expression interface {
	Node
	expressionNode()
}

// FunctionDef represents a function definition
type FunctionDef struct {
	Name       string
	Parameters []string
	Body       []Statement
}

func (f *FunctionDef) String() string {
	return fmt.Sprintf("FunctionDef{Name: %s, Params: %v, Body: %v}", f.Name, f.Parameters, f.Body)
}

func (f *FunctionDef) statementNode() {}

// FunctionCall represents a function call
type FunctionCall struct {
	Name string
	Args []Expression
}

func (f *FunctionCall) String() string {
	return fmt.Sprintf("FunctionCall{Name: %s, Args: %v}", f.Name, f.Args)
}

func (f *FunctionCall) statementNode()   {}
func (f *FunctionCall) expressionNode() {}

// SetStatement represents a variable assignment
type SetStatement struct {
	Variable string
	Value    Expression
}

func (s *SetStatement) String() string {
	return fmt.Sprintf("SetStatement{Variable: %s, Value: %s}", s.Variable, s.Value.String())
}

func (s *SetStatement) statementNode() {}

// LoopStatement represents a loop construct
type LoopStatement struct {
	Count Expression
	Body  []Statement
}

func (l *LoopStatement) String() string {
	return fmt.Sprintf("LoopStatement{Count: %s, Body: %v}", l.Count.String(), l.Body)
}

func (l *LoopStatement) statementNode() {}

// WhileStatement represents a while loop
type WhileStatement struct {
	Condition Expression
	Body      []Statement
}

func (w *WhileStatement) String() string {
	return fmt.Sprintf("WhileStatement{Condition: %s, Body: %v}", w.Condition.String(), w.Body)
}

func (w *WhileStatement) statementNode() {}

// SwitchStatement represents a switch-case construct
type SwitchStatement struct {
	Expression Expression
	Cases      []CaseClause
	Default    []Statement
}

func (s *SwitchStatement) String() string {
	return fmt.Sprintf("SwitchStatement{Expression: %s, Cases: %v, Default: %v}", s.Expression.String(), s.Cases, s.Default)
}

func (s *SwitchStatement) statementNode() {}

// CaseClause represents a case in a switch statement
type CaseClause struct {
	Values []Expression
	Body   []Statement
}

func (c *CaseClause) String() string {
	return fmt.Sprintf("CaseClause{Values: %v, Body: %v}", c.Values, c.Body)
}

// IfStatement represents a conditional statement
type IfStatement struct {
	Condition Expression
	Body      []Statement
}

func (i *IfStatement) String() string {
	return fmt.Sprintf("IfStatement{Condition: %s, Body: %v}", i.Condition.String(), i.Body)
}

func (i *IfStatement) statementNode() {}

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Value Expression
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("ReturnStatement{Value: %s}", r.Value.String())
}

func (r *ReturnStatement) statementNode() {}

// BinaryExpression represents a binary operation
type BinaryExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("BinaryExpression{Left: %s, Op: %s, Right: %s}", b.Left.String(), b.Operator, b.Right.String())
}

func (b *BinaryExpression) expressionNode() {}

// Identifier represents a variable or function name
type Identifier struct {
	Name string
}

func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier{%s}", i.Name)
}

func (i *Identifier) expressionNode() {}

// IntegerLiteral represents an integer constant
type IntegerLiteral struct {
	Value int
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("IntegerLiteral{%d}", i.Value)
}

func (i *IntegerLiteral) expressionNode() {}

// BooleanLiteral represents a boolean constant
type BooleanLiteral struct {
	Value bool
}

func (b *BooleanLiteral) String() string {
	return fmt.Sprintf("BooleanLiteral{%t}", b.Value)
}

func (b *BooleanLiteral) expressionNode() {}

// TryStatement represents a try-catch block
type TryStatement struct {
	TryBody   []Statement
	CatchBody []Statement
	ErrorVar  string
}

func (t *TryStatement) String() string {
	return fmt.Sprintf("TryStatement{Try: %v, Catch: %v, ErrorVar: %s}", t.TryBody, t.CatchBody, t.ErrorVar)
}

func (t *TryStatement) statementNode() {}

// AskStatement represents a user input prompt
type AskStatement struct {
	Question string
}

// StringLiteral represents a string constant
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral{%s}", s.Value)
}

func (s *StringLiteral) expressionNode() {}

// BracketExpression represents a parenthesized expression
type BracketExpression struct {
	Expression Expression
}

func (b *BracketExpression) String() string {
	return fmt.Sprintf("BracketExpression{%s}", b.Expression.String())
}

func (b *BracketExpression) expressionNode() {}

// AccessExpression represents object property access
type AccessExpression struct {
	Object Expression
	Key    Expression
}

func (a *AccessExpression) String() string {
	return fmt.Sprintf("AccessExpression{Object: %s, Key: %s}", a.Object.String(), a.Key.String())
}

func (a *AccessExpression) expressionNode() {}

// ArrayLiteral represents an array of expressions
type ArrayLiteral struct {
	Elements []Expression
}

func (a *ArrayLiteral) String() string {
	return fmt.Sprintf("ArrayLiteral{%v}", a.Elements)
}

func (a *ArrayLiteral) expressionNode() {}

// ObjectProperty represents a key-value pair in an object
type ObjectProperty struct {
	Key   string
	Value Expression
}

func (o *ObjectProperty) String() string {
	return fmt.Sprintf("%s: %s", o.Key, o.Value.String())
}

// ObjectLiteral represents an object with properties
type ObjectLiteral struct {
	Properties []ObjectProperty
}

func (o *ObjectLiteral) String() string {
	return fmt.Sprintf("ObjectLiteral{%v}", o.Properties)
}

func (o *ObjectLiteral) expressionNode() {}

// IndexExpression represents array or object indexing
type IndexExpression struct {
	Object Expression
	Index  Expression
}

func (i *IndexExpression) String() string {
	return fmt.Sprintf("IndexExpression{Object: %s, Index: %s}", i.Object.String(), i.Index.String())
}

func (i *IndexExpression) expressionNode() {}
