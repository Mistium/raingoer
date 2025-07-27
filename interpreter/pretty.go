package interpreter

import (
	"fmt"
	"strings"
)

func (i *Interpreter) prettyValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "nil"
	case []interface{}:
		var out []string
		for _, elem := range v {
			out = append(out, i.prettyValue(elem))
		}
		return "[" + strings.Join(out, ", ") + "]"
	case map[string]interface{}:
		var out []string
		for k, v2 := range v {
			out = append(out, fmt.Sprintf("%s: %s", k, i.prettyValue(v2)))
		}
		return "{" + strings.Join(out, ", ") + "}"
	default:
		return fmt.Sprintf("%v", v)
	}
}
