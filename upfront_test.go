package upfront

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Decodes the body or else
func decode[T any](r io.Reader) T {
	var t T
	if err := json.NewDecoder(r).Decode(&t); err != nil {
		panic(fmt.Sprintf("error decoding: %s", err))
	}

	return t
}

type Response struct {
	Name string `json:"name"`
}

type Err struct {
	Message string `json:"message"`
}

func TestOKResultHandler(t *testing.T) {
	handler := Handler[Response, Err](func(in Request) Result[Response, Err] {
		return OKResult[Response, Err](
			Response{
				Name: "James",
			},
			http.StatusOK,
		)
	})

	// Test inputs to our ServeHTTP func
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	// Should be a good code and response
	if want := http.StatusOK; resp.StatusCode != want {
		t.Fatalf("expected status code %d, instead got: %d", want, resp.StatusCode)
	}
	response := decode[Response](resp.Body)
	if response.Name != "James" {
		t.Fatalf("expected name to be James, instead got: %s", response.Name)
	}
}

func TestErrResultHandler(t *testing.T) {
	handler := Handler[Response, Err](func(in Request) Result[Response, Err] {
		return ErrResult[Response, Err](
			Err{
				Message: "could not find that record",
			},
			http.StatusNotFound,
		)
	})

	// Test inputs to our ServeHTTP func
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	// Should be a not found code and response
	if want := http.StatusNotFound; resp.StatusCode != want {
		t.Fatalf("expected status code %d, instead got: %d", want, resp.StatusCode)
	}
	e := decode[Err](resp.Body)
	if want := "could not find that record"; e.Message != want {
		t.Fatalf("expected error message to be '%s', instead got: '%s'", want, e.Message)
	}
}
