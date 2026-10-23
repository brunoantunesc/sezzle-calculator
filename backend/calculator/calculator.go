// Package calculator provides stateless float64 arithmetic independent of HTTP.
// Results follow Go floating-point semantics, including rounding, underflow to
// zero, and power(0, 0) = 1. Inputs and results must be finite; no additional
// magnitude limits are imposed.
package calculator

import (
	"errors"
	"math"
)

var (
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrInvalidOperands      = errors.New("invalid operands")
	ErrDivisionByZero       = errors.New("cannot divide by zero")
	ErrInvalidResult        = errors.New("invalid mathematical result")
)

// Calculate validates the operation and operand count and returns its result.
// Binary operations take exactly two operands; sqrt takes exactly one.
func Calculate(operation string, operands []float64) (float64, error) {
	count := 2
	switch operation {
	case "add", "subtract", "multiply", "divide", "power":
	case "sqrt":
		count = 1
	default:
		return 0, ErrUnsupportedOperation
	}
	if len(operands) != count {
		return 0, ErrInvalidOperands
	}
	for _, operand := range operands {
		if math.IsNaN(operand) || math.IsInf(operand, 0) {
			return 0, ErrInvalidOperands
		}
	}
	var result float64
	switch operation {
	case "add":
		result = operands[0] + operands[1]
	case "subtract":
		result = operands[0] - operands[1]
	case "multiply":
		result = operands[0] * operands[1]
	case "divide":
		if operands[1] == 0 {
			return 0, ErrDivisionByZero
		}
		result = operands[0] / operands[1]
	case "power":
		result = math.Pow(operands[0], operands[1])
	case "sqrt":
		result = math.Sqrt(operands[0])
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrInvalidResult
	}
	return result, nil
}
