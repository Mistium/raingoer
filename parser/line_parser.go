package parser

import (
	"strconv"
	"strings"
	"github.com/mistium/raingoer/ast"
	"github.com/mistium/raingoer/tokens"
)

type LineParser struct {
	lines []string
	pos   int
}

func NewLineParser(input string) *LineParser {
	lines := tokens.Tokenise(input, '\n')
	var filteredLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			filteredLines = append(filteredLines, line)
		}
	}
	
	return &LineParser{
		lines: filteredLines,
		pos:   0,
	}
}

func (lp *LineParser) Parse() *ast.Program {
	program := &ast.Program{}
	
	for lp.pos < len(lp.lines) {
		stmt := lp.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		lp.pos++
	}
	
	return program
}

func (lp *LineParser) parseStatement() ast.Statement {
	if lp.pos >= len(lp.lines) {
		return nil
	}
	
	line := lp.lines[lp.pos]
	tokens := lp.tokenizeLine(line)
	
	if len(tokens) == 0 { return nil }
	if strings.HasPrefix(tokens[0], "//") {
		return nil
	}

	switch tokens[0] {
	case "func": return lp.parseFunctionDef()
	case "set": return lp.parseSetStatement(tokens)
	case "loop": return lp.parseLoopStatement(tokens)
	case "while": return lp.parseWhileStatement(tokens)
	case "switch": return lp.parseSwitchStatement(tokens)
	case "if": return lp.parseIfStatement(tokens)
	case "try": return lp.parseTryStatement()
	case "return": return lp.parseReturnStatement(tokens)
	default: return lp.parseFunctionCall(tokens)
	}
}

func (lp *LineParser) tokenizeLine(line string) []string {
	if commentIndex := strings.Index(line, "//"); commentIndex != -1 { line = line[:commentIndex] }
	line = strings.TrimSpace(line)
	
	if line == "" { return []string{} }
	
	var result []string
	var currentToken strings.Builder
	inBrackets := 0
	inBraces := 0
	inString := false
	
	for _, char := range line {
		if char == '"' && !inString {
			inString = true
			currentToken.WriteRune(char)
		} else if char == '"' && inString {
			inString = false
			currentToken.WriteRune(char)
		} else if inString {
			currentToken.WriteRune(char)
		} else {
			switch char {
			case '[':
				if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				result = append(result, "[")
				inBrackets++
			case ']':
				if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				result = append(result, "]")
				inBrackets--
			case '{':
				if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				result = append(result, "{")
				inBraces++
			case '}':
				if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				result = append(result, "}")
				inBraces--
			case ',':
				if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				result = append(result, ",")
			case ':':
				if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
				result = append(result, ":")
			case ' ':
				if (inBrackets > 0 || inBraces > 0) && currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				} else if currentToken.Len() > 0 {
					result = append(result, strings.TrimSpace(currentToken.String()))
					currentToken.Reset()
				}
			default:
				currentToken.WriteRune(char)
			}
		}
	}
	
	if currentToken.Len() > 0 {
		result = append(result, strings.TrimSpace(currentToken.String()))
	}

	var filtered []string
	for _, token := range result {
		if strings.TrimSpace(token) != "" {
			filtered = append(filtered, strings.TrimSpace(token))
		}
	}
	
	return filtered
}

func (lp *LineParser) parseFunctionDef() *ast.FunctionDef {
	var body []ast.Statement
	lp.pos++
	
	line := lp.lines[lp.pos-1]
	tokens := lp.tokenizeLine(line)
	
	name := tokens[1]
	var params []string
	for i := 2; i < len(tokens); i++ {
		params = append(params, tokens[i])
	}

	for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
		stmt := lp.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		lp.pos++
	}
	
	if lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) == "end" {

	}
	
	return &ast.FunctionDef{
		Name:       name,
		Parameters: params,
		Body:       body,
	}
}

func (lp *LineParser) parseSetStatement(tokens []string) *ast.SetStatement {
	if len(tokens) < 4 || tokens[2] != "to" {
		return nil
	}
	
	variable := tokens[1]
	value := lp.parseExpressionFromTokens(tokens[3:])
	
	return &ast.SetStatement{
		Variable: variable,
		Value:    value,
	}
}

func (lp *LineParser) parseLoopStatement(tokens []string) *ast.LoopStatement {
	if len(tokens) < 2 { return nil }

	count := lp.parseExpressionFromTokens(tokens[1:2])
	
	var body []ast.Statement
	lp.pos++

	for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
		stmt := lp.parseStatement()
		if stmt != nil { body = append(body, stmt) }
		lp.pos++
	}
	
	if lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) == "end" {}
	
	return &ast.LoopStatement{
		Count: count,
		Body:  body,
	}
}

func (lp *LineParser) parseWhileStatement(tokens []string) *ast.WhileStatement {
	if len(tokens) < 2 {
		return nil
	}
	
	condition := lp.parseExpressionFromTokens(tokens[1:])
	
	var body []ast.Statement
	lp.pos++
	
	for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
		stmt := lp.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		lp.pos++
	}
	
	if lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) == "end" {
	}
	
	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (lp *LineParser) parseSwitchStatement(tokens []string) *ast.SwitchStatement {
	if len(tokens) < 2 {
		return nil
	}
	
	expression := lp.parseExpressionFromTokens(tokens[1:])
	
	var cases []ast.CaseClause
	var defaultCase []ast.Statement
	lp.pos++
	
	for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
		line := strings.TrimSpace(lp.lines[lp.pos])
		lineTokens := lp.tokenizeLine(line)
		
		if len(lineTokens) == 0 {
			lp.pos++
			continue
		}
		
		if lineTokens[0] == "case" {
			caseClause := lp.parseCaseClause(lineTokens)
			if caseClause != nil {
				cases = append(cases, *caseClause)
			}
		} else if lineTokens[0] == "default" {
			lp.pos++
			for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" && !strings.HasPrefix(strings.TrimSpace(lp.lines[lp.pos]), "case") {
				stmt := lp.parseStatement()
				if stmt != nil {
					defaultCase = append(defaultCase, stmt)
				}
				lp.pos++
			}
			lp.pos-- // Back up one since the outer loop will increment
		} else {
			lp.pos++
		}
	}
	
	return &ast.SwitchStatement{
		Expression: expression,
		Cases:      cases,
		Default:    defaultCase,
	}
}

func (lp *LineParser) parseCaseClause(tokens []string) *ast.CaseClause {
	if len(tokens) < 2 {
		return nil
	}
	
	// Parse case values (can be multiple separated by commas)
	var values []ast.Expression
	var currentValue []string
	
	for i := 1; i < len(tokens); i++ {
		if tokens[i] == "," {
			if len(currentValue) > 0 {
				expr := lp.parseExpressionFromTokens(currentValue)
				if expr != nil {
					values = append(values, expr)
				}
				currentValue = nil
			}
		} else {
			currentValue = append(currentValue, tokens[i])
		}
	}
	
	if len(currentValue) > 0 {
		expr := lp.parseExpressionFromTokens(currentValue)
		if expr != nil {
			values = append(values, expr)
		}
	}
	
	// Parse case body
	var body []ast.Statement
	lp.pos++
	
	for lp.pos < len(lp.lines) {
		line := strings.TrimSpace(lp.lines[lp.pos])
		if line == "end" || strings.HasPrefix(line, "case") || strings.HasPrefix(line, "default") {
			lp.pos-- // Back up one since the caller will increment
			break
		}
		
		stmt := lp.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		lp.pos++
	}
	
	return &ast.CaseClause{
		Values: values,
		Body:   body,
	}
}

func (lp *LineParser) parseIfStatement(tokens []string) *ast.IfStatement {
	if len(tokens) < 2 {
		return nil
	}
	
	condition := lp.parseExpressionFromTokens(tokens[1:])
	
	var body []ast.Statement
	lp.pos++
	
	for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
		stmt := lp.parseStatement()
		if stmt != nil {
			body = append(body, stmt)
		}
		lp.pos++
	}
	
	if lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) == "end" {
	}
	
	return &ast.IfStatement{
		Condition: condition,
		Body:      body,
	}
}

func (lp *LineParser) parseTryStatement() *ast.TryStatement {
	var tryBody []ast.Statement
	var catchBody []ast.Statement
	var errorVar string
	
	lp.pos++
	for lp.pos < len(lp.lines) && !strings.HasPrefix(strings.TrimSpace(lp.lines[lp.pos]), "catch") && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
		stmt := lp.parseStatement()
		if stmt != nil {
			tryBody = append(tryBody, stmt)
		}
		lp.pos++
	}
	if lp.pos < len(lp.lines) && strings.HasPrefix(strings.TrimSpace(lp.lines[lp.pos]), "catch") {
		catchLine := strings.TrimSpace(lp.lines[lp.pos])
		tokens := lp.tokenizeLine(catchLine)
		if len(tokens) > 1 {
			errorVar = tokens[1]
		}
		lp.pos++
		for lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) != "end" {
			stmt := lp.parseStatement()
			if stmt != nil {
				catchBody = append(catchBody, stmt)
			}
			lp.pos++
		}
	}
	
	if lp.pos < len(lp.lines) && strings.TrimSpace(lp.lines[lp.pos]) == "end" {
	}
	
	return &ast.TryStatement{
		TryBody:   tryBody,
		CatchBody: catchBody,
		ErrorVar:  errorVar,
	}
}

func (lp *LineParser) parseReturnStatement(tokens []string) *ast.ReturnStatement {
	if len(tokens) < 2 {
		return nil
	}
	
	value := lp.parseExpressionFromTokens(tokens[1:])
	
	return &ast.ReturnStatement{
		Value: value,
	}
}

func (lp *LineParser) parseFunctionCall(tokens []string) *ast.FunctionCall {
	if len(tokens) == 0 {
		return nil
	}
	
	name := tokens[0]
	var args []ast.Expression
	if len(tokens) > 1 {
		arg := lp.parseExpressionFromTokens(tokens[1:])
		if arg != nil {
			args = append(args, arg)
		}
	}
	
	return &ast.FunctionCall{
		Name: name,
		Args: args,
	}
}

func (lp *LineParser) parseExpressionFromTokens(tokens []string) ast.Expression {
	if len(tokens) == 0 {
		return nil
	}
	
	if len(tokens) == 1 {
		return lp.parsePrimary(tokens[0])
	}
	for i := 0; i < len(tokens)-1; i++ {
		if tokens[i+1] == "{" {
			braceEnd := -1
			braceDepth := 0
			for j := i + 1; j < len(tokens); j++ {
				if tokens[j] == "{" {
					braceDepth++
				} else if tokens[j] == "}" {
					braceDepth--
					if braceDepth == 0 {
						braceEnd = j
						break
					}
				}
			}
			
			if braceEnd > i + 2 {
				objectTokens := tokens[:i+1]
				keyTokens := tokens[i+2:braceEnd]
				
				object := lp.parseExpressionFromTokens(objectTokens)
				key := lp.parseExpressionFromTokens(keyTokens)
				
				return &ast.AccessExpression{
					Object: object,
					Key:    key,
				}
			}
		}
	}

	for i := 1; i < len(tokens); i++ {
		if isOperator(tokens[i]) {
			left := lp.parseExpressionFromTokens(tokens[:i])
			right := lp.parseExpressionFromTokens(tokens[i+1:])
			return &ast.BinaryExpression{
				Left:     left,
				Operator: tokens[i],
				Right:    right,
			}
		}
	}

	if tokens[0] == "[" {
		bracketEnd := -1
		for i := 1; i < len(tokens); i++ {
			if tokens[i] == "]" {
				bracketEnd = i
				break
			}
		}
		
		if bracketEnd > 1 {
			innerTokens := tokens[1:bracketEnd]
			if len(innerTokens) > 0 {
				name := innerTokens[0]
				var args []ast.Expression
				for i := 1; i < len(innerTokens); i++ {
					arg := lp.parsePrimary(innerTokens[i])
					if arg != nil {
						args = append(args, arg)
					}
				}
				
				functionCall := &ast.FunctionCall{
					Name: name,
					Args: args,
				}
				
				return &ast.BracketExpression{Expression: functionCall}
			}
		}
	}

	if tokens[0] == "{" {
		return lp.parseArrayOrObject(tokens)
	}
	
	return lp.parsePrimary(tokens[0])
}

func (lp *LineParser) parseArrayOrObject(tokens []string) ast.Expression {
	if len(tokens) < 2 || tokens[0] != "{" {
		return nil
	}
	
	braceEnd := -1
	for i := 1; i < len(tokens); i++ {
		if tokens[i] == "}" {
			braceEnd = i
			break
		}
	}
	
	if braceEnd == -1 {
		return nil
	}
	
	innerTokens := tokens[1:braceEnd]
	if len(innerTokens) == 0 {
		return &ast.ArrayLiteral{Elements: []ast.Expression{}}
	}
	
	hasColon := false
	for _, token := range innerTokens {
		if token == ":" {
			hasColon = true
			break
		}
	}
	
	if hasColon {
		return lp.parseObject(innerTokens)
	} else {
		return lp.parseArray(innerTokens)
	}
}

func (lp *LineParser) parseArray(tokens []string) *ast.ArrayLiteral {
	   var elements []ast.Expression
	   var currentElement []string
	   depth := 0
	   for _, token := range tokens {
			   if token == "," && depth == 0 {
					   if len(currentElement) > 0 {
							   expr := lp.parseExpressionFromTokens(currentElement)
							   if expr != nil {
									   elements = append(elements, expr)
							   }
							   currentElement = nil
					   }
			   } else {
					   if token == "{" {
							   depth++
					   } else if token == "}" {
							   depth--
					   }
					   currentElement = append(currentElement, token)
			   }
	   }
	   if len(currentElement) > 0 {
			   expr := lp.parseExpressionFromTokens(currentElement)
			   if expr != nil {
					   elements = append(elements, expr)
			   }
	   }
	   filtered := make([]ast.Expression, 0, len(elements))
	   for _, e := range elements {
			   if e != nil {
					   filtered = append(filtered, e)
			   }
	   }
	   return &ast.ArrayLiteral{Elements: filtered}
}

func (lp *LineParser) parseObject(tokens []string) *ast.ObjectLiteral {
	   var properties []ast.ObjectProperty
	   var currentProp []string
	   depth := 0
	   for _, token := range tokens {
			   if token == "," && depth == 0 {
					   if len(currentProp) > 0 {
							   prop := lp.parseObjectProperty(currentProp)
							   if prop != nil {
									   properties = append(properties, *prop)
							   }
							   currentProp = nil
					   }
			   } else {
					   if token == "{" {
							   depth++
					   } else if token == "}" {
							   depth--
					   }
					   currentProp = append(currentProp, token)
			   }
	   }
	   if len(currentProp) > 0 {
			   prop := lp.parseObjectProperty(currentProp)
			   if prop != nil {
					   properties = append(properties, *prop)
			   }
	   }
	   filtered := make([]ast.ObjectProperty, 0, len(properties))
	   for _, p := range properties {
			   if p.Value != nil {
					   filtered = append(filtered, p)
			   }
	   }
	   return &ast.ObjectLiteral{Properties: filtered}
}

func (lp *LineParser) parseObjectProperty(tokens []string) *ast.ObjectProperty {
	   colonIndex := -1
	   for i, token := range tokens {
			   if token == ":" {
					   colonIndex = i
					   break
			   }
	   }
	   if colonIndex == -1 || colonIndex == 0 || colonIndex == len(tokens)-1 {
			   return nil
	   }
	   key := tokens[0]
	   valueTokens := tokens[colonIndex+1:]
	   // If value starts with {, group tokens until matching },
	   // If value starts with [, group tokens until matching ]
	   if len(valueTokens) > 0 && (valueTokens[0] == "{" || valueTokens[0] == "[") {
			   open, close := valueTokens[0], "}"
			   if open == "[" {
					   close = "]"
			   }
			   depth := 0
			   end := 0
			   for i, t := range valueTokens {
					   if t == open {
							   depth++
					   } else if t == close {
							   depth--
							   if depth == 0 {
									   end = i + 1
									   break
							   }
					   }
			   }
			   valueTokens = valueTokens[:end]
	   }
	   var value ast.Expression
	   if len(valueTokens) == 1 {
			   token := valueTokens[0]
			   if val, err := strconv.Atoi(token); err == nil {
					   value = &ast.IntegerLiteral{Value: val}
			   } else if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
					   value = &ast.StringLiteral{Value: strings.Trim(token, "\"")}
			   } else {
					   value = &ast.StringLiteral{Value: token}
			   }
	   } else {
			   value = lp.parseExpressionFromTokens(valueTokens)
	   }
	   if value == nil {
			   return nil
	   }
	   return &ast.ObjectProperty{
			   Key:   key,
			   Value: value,
	   }
}

func (lp *LineParser) parsePrimary(token string) ast.Expression {
	if val, err := strconv.Atoi(token); err == nil {
		return &ast.IntegerLiteral{Value: val}
	}

	if token == "true" {
		return &ast.BooleanLiteral{Value: true}
	}

	if token == "false" {
		return &ast.BooleanLiteral{Value: false}
	}

	if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return &ast.StringLiteral{Value: strings.Trim(token, "\"")}
	}

	return &ast.Identifier{Name: token}
}
