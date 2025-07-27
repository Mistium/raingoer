package tokens

func tokenise(code string, delimiter byte) []string {
	var (
		out     []byte
		split   []string
		escaped bool
		brackets int
		bDepth   int
	)

	for i := 0; i < len(code); i++ {
		ch := code[i]

		if brackets == 0 && !escaped {
			if ch == '[' || ch == '{' || ch == '(' {
				bDepth++
			}
			if ch == ']' || ch == '}' || ch == ')' {
				bDepth--
				if bDepth < 0 {
					bDepth = 0
				}
			}
		}

		if ch == '"' && !escaped {
			brackets = 1 - brackets
			out = append(out, '"')
		} else if ch == '\\' && !escaped {
			escaped = true
			out = append(out, '\\')
			continue
		} else {
			out = append(out, ch)
			escaped = false
		}

		if i+1 < len(code) && code[i+1] == delimiter && brackets == 0 && bDepth == 0 {
			split = append(split, string(out))
			out = nil
			i++
		}
	}

	if len(out) > 0 {
		split = append(split, string(out))
	}
	return split
}
