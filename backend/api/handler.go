// Package api exposes the calculator over HTTP.
package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"github.com/brunoantunesc/sezzle-calculator/backend/calculator"
)

// NewHandler returns the calculator and health routes.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, "GET") {
			return
		}
		writeJSON(w, 200, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/api/v1/calculations", calculate)
	return mux
}

func requireMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method == method {
		return true
	}
	w.Header().Set("Allow", method)
	writeError(w, 405, "METHOD_NOT_ALLOWED", "Method not allowed")
	return false
}

func calculate(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, "POST") {
		return
	}
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		writeError(w, 415, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json")
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 4096))
	if err != nil {
		var limit *http.MaxBytesError
		if errors.As(err, &limit) {
			writeError(w, 413, "PAYLOAD_TOO_LARGE", "Request body exceeds 4096 bytes")
		} else {
			writeError(w, 400, "INVALID_REQUEST", "Cannot read request body")
		}
		return
	}
	// Decode fields individually to reject duplicate keys and enforce exact names.
	// Validating the whole document first also rejects trailing JSON data.
	if !json.Valid(body) {
		writeError(w, 400, "INVALID_REQUEST", "Invalid JSON request")
		return
	}
	fields, ok := requestFields(body)
	if !ok {
		writeError(w, 400, "INVALID_REQUEST", "Expected an object with operation and operands fields")
		return
	}
	var operation *string
	if json.Unmarshal(fields["operation"], &operation) != nil || operation == nil {
		writeError(w, 400, "INVALID_REQUEST", "Operation must be a string")
		return
	}
	var values []json.RawMessage
	if json.Unmarshal(fields["operands"], &values) != nil || values == nil {
		writeError(w, 400, "INVALID_OPERANDS", "Operands must be an array of numbers")
		return
	}
	count := 2
	switch *operation {
	case "sqrt":
		count = 1
	case "add", "subtract", "multiply", "divide", "power":
	default:
		writeError(w, 400, "UNSUPPORTED_OPERATION", "Unsupported operation")
		return
	}
	if len(values) != count {
		writeError(w, 400, "INVALID_OPERANDS", "Invalid operand count or value")
		return
	}
	operands := make([]float64, len(values))
	for i, v := range values {
		var number *float64
		if json.Unmarshal(v, &number) != nil || number == nil {
			writeError(w, 400, "INVALID_OPERANDS", "Operands must be numbers")
			return
		}
		operands[i] = *number
	}
	result, err := calculator.Calculate(*operation, operands)
	switch {
	case err == nil:
		writeJSON(w, 200, struct {
			Result float64 `json:"result"`
		}{result})
	case errors.Is(err, calculator.ErrUnsupportedOperation):
		writeError(w, 400, "UNSUPPORTED_OPERATION", "Unsupported operation")
	case errors.Is(err, calculator.ErrInvalidOperands):
		writeError(w, 400, "INVALID_OPERANDS", "Invalid operand count or value")
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, 422, "DIVISION_BY_ZERO", "Cannot divide by zero")
	case errors.Is(err, calculator.ErrInvalidResult):
		writeError(w, 422, "INVALID_RESULT", "Invalid mathematical result")
	default:
		writeError(w, 500, "INTERNAL_ERROR", "Internal server error")
	}
}

func requestFields(body []byte) (map[string]json.RawMessage, bool) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		return nil, false
	}
	fields := make(map[string]json.RawMessage)
	for decoder.More() {
		token, err = decoder.Token()
		if err != nil {
			return nil, false
		}
		key, ok := token.(string)
		if !ok || (key != "operation" && key != "operands") {
			return nil, false
		}
		if _, exists := fields[key]; exists {
			return nil, false
		}
		var value json.RawMessage
		if decoder.Decode(&value) != nil {
			return nil, false
		}
		fields[key] = value
	}
	token, err = decoder.Token()
	if err != nil || token != json.Delim('}') {
		return nil, false
	}
	if _, err = decoder.Token(); err != io.EOF {
		return nil, false
	}
	return fields, true
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// These responses contain only known strings and finite calculator results.
	// A write failure means the client connection is no longer usable.
	_ = json.NewEncoder(w).Encode(value)
}
