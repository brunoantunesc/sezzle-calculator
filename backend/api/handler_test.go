package api

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	valid := `{"operation":"divide","operands":[10,2]}`
	tests := []struct {
		name, method, path, content, body string
		status                            int
		code                              string
	}{
		{"calculation", "POST", "/api/v1/calculations", "application/json; charset=utf-8", valid, 200, ""},
		{"health", "GET", "/health", "", "", 200, ""},
		{"method", "GET", "/api/v1/calculations", "", "", 405, "METHOD_NOT_ALLOWED"},
		{"health method", "POST", "/health", "", "", 405, "METHOD_NOT_ALLOWED"},
		{"head", "HEAD", "/health", "", "", 405, "METHOD_NOT_ALLOWED"},
		{"content", "POST", "/api/v1/calculations", "text/plain", valid, 415, "UNSUPPORTED_MEDIA_TYPE"},
		{"missing content", "POST", "/api/v1/calculations", "", valid, 415, "UNSUPPORTED_MEDIA_TYPE"},
	}
	for _, tc := range []struct {
		name, body, code string
		status           int
	}{
		{"empty", "", "INVALID_REQUEST", 400},
		{"empty operation", `{"operation":"","operands":[1,2]}`, "UNSUPPORTED_OPERATION", 400},
		{"empty operands", `{"operation":"add","operands":[]}`, "INVALID_OPERANDS", 400},
		{"one byte above limit", valid + strings.Repeat(" ", 4097-len(valid)), "PAYLOAD_TOO_LARGE", 413},
		{"broken", `{"operation":`, "INVALID_REQUEST", 400},
		{"null", `null`, "INVALID_REQUEST", 400},
		{"array", `[]`, "INVALID_REQUEST", 400},
		{"unknown", `{"operation":"add","operands":[1,2],"extra":1}`, "INVALID_REQUEST", 400},
		{"case", `{"Operation":"add","operands":[1,2]}`, "INVALID_REQUEST", 400},
		{"duplicate", `{"operation":"add","operation":"divide","operands":[1,2]}`, "INVALID_REQUEST", 400},
		{"trailing", valid + ` {}`, "INVALID_REQUEST", 400},
		{"trailing garbage", valid + ` x`, "INVALID_REQUEST", 400},
		{"missing operation", `{"operands":[1,2]}`, "INVALID_REQUEST", 400},
		{"typed operation", `{"operation":1,"operands":[1,2]}`, "INVALID_REQUEST", 400},
		{"null operation", `{"operation":null,"operands":[1,2]}`, "INVALID_REQUEST", 400},
		{"unsupported", `{"operation":"nope","operands":[1,2]}`, "UNSUPPORTED_OPERATION", 400},
		{"missing operands", `{"operation":"add"}`, "INVALID_OPERANDS", 400},
		{"null operands", `{"operation":"add","operands":null}`, "INVALID_OPERANDS", 400},
		{"typed operands", `{"operation":"add","operands":"1,2"}`, "INVALID_OPERANDS", 400},
		{"null element", `{"operation":"add","operands":[null,2]}`, "INVALID_OPERANDS", 400},
		{"string element", `{"operation":"add","operands":["1",2]}`, "INVALID_OPERANDS", 400},
		{"out of range", `{"operation":"add","operands":[1e400,2]}`, "INVALID_OPERANDS", 400},
		{"wrong count", `{"operation":"sqrt","operands":[1,2]}`, "INVALID_OPERANDS", 400},
		{"zero divisor", `{"operation":"divide","operands":[1,0]}`, "DIVISION_BY_ZERO", 422},
		{"negative root", `{"operation":"sqrt","operands":[-1]}`, "INVALID_RESULT", 422},
		{"overflow", `{"operation":"power","operands":[2,1024]}`, "INVALID_RESULT", 422},
		{"oversized", valid + strings.Repeat(" ", 4096), "PAYLOAD_TOO_LARGE", 413},
		{"exact limit", valid + strings.Repeat(" ", 4096-len(valid)), "", 200},
	} {
		tests = append(tests, struct {
			name, method, path, content, body string
			status                            int
			code                              string
		}{tc.name, "POST", "/api/v1/calculations", "application/json", tc.body, tc.status, tc.code})
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", tc.content)
			w := httptest.NewRecorder()
			NewHandler().ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Fatalf("status %d, want %d: %s", w.Code, tc.status, w.Body.String())
			}
			if w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("missing JSON content type")
			}
			var response struct {
				Result *float64                        `json:"result"`
				Status string                          `json:"status"`
				Error  *struct{ Code, Message string } `json:"error"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if tc.code != "" {
				if response.Error == nil || response.Error.Code != tc.code || response.Error.Message == "" {
					t.Fatalf("unexpected error: %s", w.Body.String())
				}
			} else if tc.path == "/health" {
				if response.Status != "ok" {
					t.Fatal("invalid health")
				}
			} else if response.Result == nil || *response.Result != 5 {
				t.Fatal("invalid result")
			}
			if tc.status == 405 && w.Header().Get("Allow") == "" {
				t.Fatal("missing Allow")
			}
		})
	}
}

func TestCalculationOperations(t *testing.T) {
	for _, tc := range []struct {
		body string
		want float64
	}{
		{`{"operation":"add","operands":[-2,3.5]}`, 1.5},
		{`{"operation":"subtract","operands":[0,4]}`, -4},
		{`{"operation":"multiply","operands":[-2,-3]}`, 6},
		{`{"operation":"power","operands":[2,-2]}`, 0.25},
		{`{"operation":"sqrt","operands":[2.25]}`, 1.5},
	} {
		req := httptest.NewRequest("POST", "/api/v1/calculations", strings.NewReader(tc.body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		NewHandler().ServeHTTP(w, req)
		var response struct{ Result float64 }
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if w.Code != 200 || response.Result != tc.want {
			t.Errorf("request %s: status %d, body %s", tc.body, w.Code, w.Body.String())
		}
	}
}

func TestRequestFieldsDocument(t *testing.T) {
	for _, body := range []string{`{"operation":"add"`, `{"operation":"add"]`, `{} {}`, `{} x`} {
		if _, ok := requestFields([]byte(body)); ok {
			t.Errorf("accepted invalid document %q", body)
		}
	}
	if _, ok := requestFields([]byte(`{} `)); !ok {
		t.Fatal("rejected closed object with whitespace")
	}
}
