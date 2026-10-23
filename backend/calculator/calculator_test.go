package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name, operation string
		operands        []float64
		want            float64
		err             error
	}{
		{"addition", "add", []float64{-2, 3.5}, 1.5, nil},
		{"subtraction", "subtract", []float64{0, 4}, -4, nil},
		{"multiplication", "multiply", []float64{-2, -3}, 6, nil},
		{"division", "divide", []float64{5, 2}, 2.5, nil},
		{"power", "power", []float64{2, -2}, 0.25, nil},
		{"root", "sqrt", []float64{2.25}, 1.5, nil},
		{"zero root", "sqrt", []float64{0}, 0, nil},
		{"zero power", "power", []float64{0, 0}, 1, nil},
		{"unsupported", "ADD", []float64{1, 2}, 0, ErrUnsupportedOperation},
		{"missing", "add", nil, 0, ErrInvalidOperands},
		{"binary too few", "divide", []float64{1}, 0, ErrInvalidOperands},
		{"binary too many", "multiply", []float64{1, 2, 3}, 0, ErrInvalidOperands},
		{"unary too many", "sqrt", []float64{1, 2}, 0, ErrInvalidOperands},
		{"unary missing", "sqrt", nil, 0, ErrInvalidOperands},
		{"zero divisor", "divide", []float64{1, 0}, 0, ErrDivisionByZero},
		{"negative zero divisor", "divide", []float64{1, math.Copysign(0, -1)}, 0, ErrDivisionByZero},
		{"negative root", "sqrt", []float64{-1}, 0, ErrInvalidResult},
		{"fractional negative power", "power", []float64{-1, 0.5}, 0, ErrInvalidResult},
		{"power overflow", "power", []float64{2, 1024}, 0, ErrInvalidResult},
		{"addition overflow", "add", []float64{math.MaxFloat64, math.MaxFloat64}, 0, ErrInvalidResult},
		{"negative overflow", "multiply", []float64{-math.MaxFloat64, 2}, 0, ErrInvalidResult},
		{"nan operand", "add", []float64{math.NaN(), 1}, 0, ErrInvalidOperands},
		{"infinite operand", "power", []float64{math.Inf(1), 0}, 0, ErrInvalidOperands},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.operation, tt.operands)
			if !errors.Is(err, tt.err) {
				t.Fatalf("error = %v, want %v", err, tt.err)
			}
			if err == nil && got != tt.want {
				t.Errorf("result = %v, want %v", got, tt.want)
			}
		})
	}
}
