# Raingoer Language Documentation

## Overview

Raingoer is a simple interpreted language made because i was bored lmfao

## Syntax Reference

### Function Definition

```go
func add x y
  set result to x + y
  return result
end
```

### Variable Assignment

```go
set x to 42
set name to "Alice"
set sum to a + b
```

### Loops

```go
loop 5
  state [add 1 2]
end
```

### Return Statement

```go
return x * 2
```

### Function Call

```go
add 5 3
```

### Bracket Evaluation

Use brackets to evaluate a function call as an expression:

```go
set result to [add 2 3]
state [add 10 20]
```

### Built-in Functions

- `state [expression]` — Prints the result of evaluating the expression.

  Example:

  ```go
  state [add 5 3]
  state [fib 10]
  state [ask "how are you? "]
  ```

## Operators

### Arithmetic Operators

| Operator | Description     | Example         | Result |
|----------|----------------|-----------------|--------|
| +        | Addition       | 1 + 2           | 3      |
| -        | Subtraction    | 5 - 3           | 2      |
| *        | Multiplication | 4 * 2           | 8      |
| /        | Division       | 8 / 2           | 4      |

### Comparison Operators

| Operator | Description     | Example         | Result |
|----------|----------------|-----------------|--------|
| ==       | Equal          | 2 == 2          | true   |
| !=       | Not equal      | 2 != 3          | true   |
| <        | Less than      | 1 < 2           | true   |
| >        | Greater than   | 3 > 2           | true   |
| <=       | Less or equal  | 2 <= 2          | true   |
| >=       | Greater/equal  | 3 >= 2          | true   |

## Example: Fibonacci Benchmark

```go
func fib v
  set a to 0
  set b to 1
  set c to 0
  loop v
    set c to a + b
    set a to b
    set b to c
  end
  return c
end

loop 1000
  fib 1000
end
```

## Example: Function Composition

```go
func add x y
  set result to x + y
  return result
end

state [add 5 3]
```

## Features

- Functions with parameters and return values
- Variable assignment and arithmetic
- Loops
- Bracket-based function evaluation
- Built-in `state` for output

## License

MIT
