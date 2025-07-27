package interpreter

import (
	"fmt"
	"strings"
	"github.com/mistium/raingoer/ast"
)

type Environment struct {
	variables map[string]interface{}
	functions map[string]*ast.FunctionDef
	parent    *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		variables: make(map[string]interface{}),
		functions: make(map[string]*ast.FunctionDef),
		parent:    parent,
	}
}

func (env *Environment) Get(name string) (interface{}, bool) {
	if val, ok := env.variables[name]; ok {
		return val, true
	}
	if env.parent != nil {
		return env.parent.Get(name)
	}
	return nil, false
}

func (env *Environment) Set(name string, value interface{}) {
	if _, ok := env.variables[name]; ok {
		env.variables[name] = value
		return
	}
	if env.parent != nil {
		if _, ok := env.parent.Get(name); ok {
			env.parent.Set(name, value)
			return
		}
	}
	env.variables[name] = value
}

func (env *Environment) GetFunction(name string) (*ast.FunctionDef, bool) {
	if fn, ok := env.functions[name]; ok {
		return fn, true
	}
	if env.parent != nil {
		return env.parent.GetFunction(name)
	}
	return nil, false
}

func (env *Environment) SetFunction(name string, fn *ast.FunctionDef) {
	env.functions[name] = fn
}

type Interpreter struct {
	env *Environment
}

func New() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(nil),
	}
}

func (i *Interpreter) Interpret(program *ast.Program) interface{} {
	var result interface{}
	
	for _, stmt := range program.Statements {
		currentResult := i.evalStatement(stmt)
		if currentResult != nil {
			result = currentResult
		}
	}
	
	return result
}

func (i *Interpreter) evalStatement(stmt ast.Statement) interface{} {
	switch node := stmt.(type) {
	case *ast.FunctionDef:
		i.env.SetFunction(node.Name, node)
		return nil
		
	case *ast.FunctionCall:
		return i.evalFunctionCall(node)
		
	case *ast.SetStatement:
		value := i.evalExpression(node.Value)
		i.env.Set(node.Variable, value)
		return nil
		
	case *ast.LoopStatement:
		count := i.evalExpression(node.Count)
		if countInt, ok := count.(int); ok {
			loopEnv := NewEnvironment(i.env)
			oldEnv := i.env
			i.env = loopEnv
			
			for j := 0; j < countInt; j++ {
				for _, bodyStmt := range node.Body {
					result := i.evalStatement(bodyStmt)
					if _, isReturn := bodyStmt.(*ast.ReturnStatement); isReturn {
						i.env = oldEnv
						return result
					}
				}
			}
			
			i.env = oldEnv
		}
		return nil
		
	case *ast.WhileStatement:
		whileEnv := NewEnvironment(i.env)
		oldEnv := i.env
		i.env = whileEnv
		
		for {
			condition := i.evalExpression(node.Condition)
			condBool, ok := condition.(bool)
			if !ok || !condBool {
				break
			}
			
			for _, bodyStmt := range node.Body {
				result := i.evalStatement(bodyStmt)
				if _, isReturn := bodyStmt.(*ast.ReturnStatement); isReturn {
					i.env = oldEnv
					return result
				}
			}
		}
		
		i.env = oldEnv
		return nil
		
	case *ast.SwitchStatement:
		switchValue := i.evalExpression(node.Expression)
		
		switchEnv := NewEnvironment(i.env)
		oldEnv := i.env
		i.env = switchEnv
		
		matched := false
		for _, caseClause := range node.Cases {
			for _, caseValue := range caseClause.Values {
				caseVal := i.evalExpression(caseValue)
				if i.valuesEqual(switchValue, caseVal) {
					matched = true
					for _, bodyStmt := range caseClause.Body {
						result := i.evalStatement(bodyStmt)
						if _, isReturn := bodyStmt.(*ast.ReturnStatement); isReturn {
							i.env = oldEnv
							return result
						}
					}
					break
				}
			}
			if matched {
				break
			}
		}
		
		// Execute default case if no case matched
		if !matched && len(node.Default) > 0 {
			for _, bodyStmt := range node.Default {
				result := i.evalStatement(bodyStmt)
				if _, isReturn := bodyStmt.(*ast.ReturnStatement); isReturn {
					i.env = oldEnv
					return result
				}
			}
		}
		
		i.env = oldEnv
		return nil
		
	case *ast.IfStatement:
		condition := i.evalExpression(node.Condition)
		if condBool, ok := condition.(bool); ok && condBool {
			ifEnv := NewEnvironment(i.env)
			oldEnv := i.env
			i.env = ifEnv
			
			for _, bodyStmt := range node.Body {
				result := i.evalStatement(bodyStmt)
				if _, isReturn := bodyStmt.(*ast.ReturnStatement); isReturn {
					i.env = oldEnv
					return result
				}
			}
			
			i.env = oldEnv
		}
		return nil
		
	case *ast.ReturnStatement:
		return i.evalExpression(node.Value)
		
	case *ast.TryStatement:
		return i.evalTryStatement(node)
		
	default:
		panic(fmt.Sprintf("unknown statement type: %T", stmt))
	}
}

func (i *Interpreter) valuesEqual(a, b interface{}) bool {
	return a == b
}

func (i *Interpreter) evalFunctionCall(call *ast.FunctionCall) interface{} {
	   if call.Name == "state" {
			   if len(call.Args) > 0 {
					   result := i.evalExpression(call.Args[0])
					   fmt.Printf("%s\n", i.prettyValue(result))
					   return result
			   }
			   return nil
	   }



	   if call.Name == "ask" {
			   var prompt string
			   if len(call.Args) > 0 {
					   promptVal := i.evalExpression(call.Args[0])
					   prompt, _ = promptVal.(string)
			   }
			   if prompt != "" {
					   fmt.Print(prompt)
			   }
			   var input string
			   fmt.Scanln(&input)
			   return input
	   }
	
	fn, exists := i.env.GetFunction(call.Name)
	if !exists {
		panic(fmt.Sprintf("Error: Function '%s' is not defined", call.Name))
	}
	globalEnv := i.env
	for globalEnv.parent != nil {
		globalEnv = globalEnv.parent
	}
	funcEnv := NewEnvironment(globalEnv)
	oldEnv := i.env
	i.env = funcEnv
	
	for idx, param := range fn.Parameters {
		if idx < len(call.Args) {
			value := i.evalExpression(call.Args[idx])
			i.env.Set(param, value)
		}
	}
	
	var result interface{}
	for _, stmt := range fn.Body {
		result = i.evalStatement(stmt)
		if _, isReturn := stmt.(*ast.ReturnStatement); isReturn {
			break
		}
	}
	
	i.env = oldEnv
	return result
}

func (i *Interpreter) evalTryStatement(stmt *ast.TryStatement) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if len(stmt.CatchBody) > 0 {
				catchEnv := NewEnvironment(i.env)
				oldEnv := i.env
				i.env = catchEnv
				if stmt.ErrorVar != "" {
					errorMsg := fmt.Sprintf("%v", r)
					i.env.Set(stmt.ErrorVar, errorMsg)
				}
				for _, catchStmt := range stmt.CatchBody {
					i.evalStatement(catchStmt)
				}
				
				i.env = oldEnv
			}
		}
	}()
	for _, tryStmt := range stmt.TryBody {
		result := i.evalStatement(tryStmt)
		if _, isReturn := tryStmt.(*ast.ReturnStatement); isReturn {
			return result
		}
	}
	
	return nil
}

func (i *Interpreter) evalExpression(expr ast.Expression) interface{} {
	switch node := expr.(type) {
	case *ast.IntegerLiteral:
		return node.Value
		
	case *ast.BooleanLiteral:
		return node.Value
		
	case *ast.StringLiteral:
		return node.Value
		
	case *ast.Identifier:
		if val, ok := i.env.Get(node.Name); ok {
			return val
		}
		panic(fmt.Sprintf("Error: Variable '%s' is not defined in the current scope", node.Name))
		

	   case *ast.BinaryExpression:
			   left := i.evalExpression(node.Left)
			   right := i.evalExpression(node.Right)
			   if leftInt, ok1 := left.(int); ok1 {
					   if rightInt, ok2 := right.(int); ok2 {
							   switch node.Operator {
							   case "+":
									   return leftInt + rightInt
							   case "-":
									   return leftInt - rightInt
							   case "*":
									   return leftInt * rightInt
							   case "/":
									   if rightInt == 0 {
											   panic("division by zero")
									   }
									   return leftInt / rightInt
							   case "%":
									   if rightInt == 0 {
											   panic("modulo by zero")
									   }
									   return leftInt % rightInt
							   case "==":
									   return leftInt == rightInt
							   case "!=":
									   return leftInt != rightInt
							   case "<":
									   return leftInt < rightInt
							   case ">":
									   return leftInt > rightInt
							   case "<=":
									   return leftInt <= rightInt
							   case ">=":
									   return leftInt >= rightInt
							   }
					   }
			   }
			   leftStr, leftIsStr := left.(string)
			   rightStr, rightIsStr := right.(string)
			   switch node.Operator {
			   case "++":
					   if leftIsStr && rightIsStr {
							   return leftStr + rightStr
					   }
					   if leftIsStr {
							   return leftStr + fmt.Sprintf("%v", right)
					   }
					   if rightIsStr {
							   return fmt.Sprintf("%v", left) + rightStr
					   }
					   return fmt.Sprintf("%v%v", left, right)
			   case "==":
					   return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
			   case "!=":
					   return fmt.Sprintf("%v", left) != fmt.Sprintf("%v", right)
			   }

			   panic(fmt.Sprintf("cannot apply operator %s to %T and %T", node.Operator, left, right))
		
	case *ast.FunctionCall:
		return i.evalFunctionCall(node)
		
	case *ast.BracketExpression:
		return i.evalExpression(node.Expression)
		
	case *ast.ArrayLiteral:
			   var elements []interface{}
			   for _, elem := range node.Elements {
					   if elem == nil {
							   continue
					   }
					   value := i.evalExpression(elem)
					   if value != nil {
							   elements = append(elements, value)
					   }
			   }
			   return elements
		
	case *ast.ObjectLiteral:
			   obj := make(map[string]interface{})
			   for _, prop := range node.Properties {
					   if prop.Value == nil {
							   continue
					   }
					   value := i.evalExpression(prop.Value)
					   key := prop.Key
					   if strings.HasPrefix(key, "\"") && strings.HasSuffix(key, "\"") {
							   key = strings.Trim(key, "\"")
					   }
					   obj[key] = value
			   }
			   return obj
		
	case *ast.IndexExpression:
			   object := i.evalExpression(node.Object)
			   index := i.evalExpression(node.Index)
			   if object == nil {
					   return nil
			   }
			   if arr, ok := object.([]interface{}); ok {
					   if idx, ok := index.(int); ok {
							   if idx >= 0 && idx < len(arr) {
									   return arr[idx]
							   }
							   panic(fmt.Sprintf("array index %d out of bounds", idx))
					   }
					   panic("array index must be an integer")
			   }
			   if obj, ok := object.(map[string]interface{}); ok {
					   if key, ok := index.(string); ok {
							   if val, exists := obj[key]; exists {
									   return val
							   }
							   return nil
					   }
					   panic("object key must be a string")
			   }
			   return nil
		
	case *ast.AccessExpression:
			   object := i.evalExpression(node.Object)
			   key := i.evalExpression(node.Key)
			   if object == nil {
					   return nil
			   }
			   if arr, ok := object.([]interface{}); ok {
					   if idx, ok := key.(int); ok {
							   if idx >= 0 && idx < len(arr) {
									   return arr[idx]
							   }
							   panic(fmt.Sprintf("array index %d out of bounds", idx))
					   }
					   panic("array index must be an integer")
			   }
			   if obj, ok := object.(map[string]interface{}); ok {
					   if keyStr, ok := key.(string); ok {
							   if val, exists := obj[keyStr]; exists {
									   return val
							   }
							   return nil
					   }
					   panic("object key must be a string")
			   }
			   return nil
		
	   default:
			   // Return nil for unknown expression types to avoid panic and help debug
			   return nil
	}
}
