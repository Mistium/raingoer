package parser

import (
	"strconv"
	"strings"
	"github.com/mistium/raingoer/ast"
	"github.com/mistium/raingoer/tokens"
)

type Parser struct {
	tokens []string
	pos    int
}

func (p *Parser) DebugTokens() []string {
	return p.tokens
}

func New(input string) *Parser {
	lines := tokens.Tokenise(input, '\n')
	var allTokens []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var currentToken strings.Builder
		inBrackets := 0
		
		for i, char := range line {
			switch char {
			case '[':
				if currentToken.Len() > 0 {
					allTokens = append(allTokens, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				allTokens = append(allTokens, "[")
				inBrackets++
			case ']':
				if currentToken.Len() > 0 {
					allTokens = append(allTokens, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				allTokens = append(allTokens, "]")
				inBrackets--
			case ' ':
				if inBrackets > 0 || currentToken.Len() > 0 {
					if currentToken.Len() > 0 {
						allTokens = append(allTokens, strings.TrimSpace(currentToken.String()))
						currentToken.Reset()
					}
				}
			default:
				currentToken.WriteRune(char)
			}

			if i == len(line)-1 && currentToken.Len() > 0 {
				allTokens = append(allTokens, strings.TrimSpace(currentToken.String()))
			}
		}

		allTokens = append(allTokens, "\n")
	}

	var filteredTokens []string
	for _, token := range allTokens {
		if strings.TrimSpace(token) != "" {
			filteredTokens = append(filteredTokens, token)
		}
	}
	
	return &Parser{
		tokens: filteredTokens,
		pos:    0,
	}
}

func (p *Parser) currentToken() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

func (p *Parser) peekToken() string {
	if p.pos+1 >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) nextToken() {
	p.pos++
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{}
	
	for p.pos < len(p.tokens) {

		if p.currentToken() == "\n" {
			p.nextToken()
			continue
		}
		
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	startPos := p.pos
	
	switch p.currentToken() {
	case "func":
		return p.parseFunctionDef()
	case "set":
		return p.parseSetStatement()
	case "loop":
		return p.parseLoopStatement()
	case "return":
		return p.parseReturnStatement()
	default:

		if p.peekToken() != "" && !isOperator(p.peekToken()) {
			return p.parseFunctionCall()
		}

		p.nextToken()

		if p.pos == startPos {
			p.nextToken()
		}
		return nil
	}
}

func (p *Parser) parseFunctionDef() *ast.FunctionDef {
	p.nextToken()
	
	name := p.currentToken()
	p.nextToken()
	
	var params []string

	for p.currentToken() != "" && !isKeyword(p.currentToken()) {
		params = append(params, p.currentToken())
		p.nextToken()
	}
	
	var body []ast.Statement
	for p.currentToken() != "end" && p.pos < len(p.tokens) {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}
	
	if p.currentToken() == "end" {
		p.nextToken()
	}
	
	return &ast.FunctionDef{
		Name:       name,
		Parameters: params,
		Body:       body,
	}
}

func (p *Parser) parseSetStatement() *ast.SetStatement {
	p.nextToken()
	
	variable := p.currentToken()
	p.nextToken()
	
	if p.currentToken() == "to" {
		p.nextToken()
	}
	
	value := p.parseExpression()
	
	return &ast.SetStatement{
		Variable: variable,
		Value:    value,
	}
}

func (p *Parser) parseLoopStatement() *ast.LoopStatement {
	p.nextToken()
	
	count := p.parseExpression()
	
	var body []ast.Statement
	for p.currentToken() != "end" && p.pos < len(p.tokens) {
		stmt := p.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
	}
	
	if p.currentToken() == "end" {
		p.nextToken()
	}
	
	return &ast.LoopStatement{
		Count: count,
		Body:  body,
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.nextToken()
	
	value := p.parseExpression()
	
	return &ast.ReturnStatement{
		Value: value,
	}
}

func (p *Parser) parseFunctionCall() *ast.FunctionCall {
	name := p.currentToken()
	p.nextToken()
	
	var args []ast.Expression

	for p.currentToken() != "" && !isKeyword(p.currentToken()) && p.currentToken() != "\n" {
		arg := p.parseExpression()
		if arg != nil {
			args = append(args, arg)
		}

		if p.pos >= len(p.tokens) || p.currentToken() == "\n" {
			break
		}
	}
	
	return &ast.FunctionCall{
		Name: name,
		Args: args,
	}
}

func (p *Parser) parseExpression() ast.Expression {
	left := p.parsePrimary()

	if isOperator(p.currentToken()) {
		operator := p.currentToken()
		p.nextToken()
		right := p.parseExpression()
		
		return &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	
	return left
}

func (p *Parser) parsePrimary() ast.Expression {
	token := p.currentToken()

	if token == "[" {
		p.nextToken()

		if p.currentToken() != "]" {
			funcName := p.currentToken()
			p.nextToken()
			
			var args []ast.Expression

			for p.currentToken() != "]" && p.pos < len(p.tokens) {
				arg := p.parsePrimary()
				if arg != nil {
					args = append(args, arg)
				}
			}
			
			if p.currentToken() == "]" {
				p.nextToken()
			}
			
			functionCall := &ast.FunctionCall{
				Name: funcName,
				Args: args,
			}
			
			return &ast.BracketExpression{Expression: functionCall}
		}
		
		if p.currentToken() == "]" {
			p.nextToken()
		}
		return &ast.BracketExpression{Expression: &ast.IntegerLiteral{Value: 0}}
	}
	
	p.nextToken()

	if val, err := strconv.Atoi(token); err == nil {
		return &ast.IntegerLiteral{Value: val}
	}

	if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return &ast.StringLiteral{Value: strings.Trim(token, "\"")}
	}

	return &ast.Identifier{Name: token}
}

func isKeyword(token string) bool {
	keywords := []string{"func", "end", "set", "to", "loop", "return"}
	for _, keyword := range keywords {
		if token == keyword {
			return true
		}
	}
	return false
}

func isOperator(token string) bool {
	operators := []string{"+", "-", "*", "/", "==", "!=", "<", ">", "<=", ">=", "++"}
	for _, op := range operators {
		if token == op {
			return true
		}
	}
	return false
}
