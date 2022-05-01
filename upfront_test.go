package upfront

import (
	"bytes"
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

// Encodes the body or else
func encode[T any](t T) *bytes.Buffer {
	byts, err := json.Marshal(t)
	if err != nil {
		panic(fmt.Sprintf("error encoding: %s", err))
	}

	return bytes.NewBuffer(byts)
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
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("expected status code %d, instead got: %d", want, got)
	}
	response := decode[Response](resp.Body)
	if want, got := "James", response.Name; want != got {
		t.Fatalf("expected name to be '%s', instead got: '%s'", want, got)
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
	if want, got := http.StatusNotFound, resp.StatusCode; want != got {
		t.Fatalf("expected status code %d, instead got: %d", want, resp.StatusCode)
	}
	e := decode[Err](resp.Body)
	if want, got := "could not find that record", e.Message; want != got {
		t.Fatalf("expected error message to be '%s', instead got: '%s'", want, got)
	}
}

type Input struct {
	Name string `json:"name"`
}

func TestOKResultBodyHandler(t *testing.T) {
	handler := BodyHandler[Input, Response, Err](func(in BodyRequest[Input]) Result[Response, Err] {
		if want, got := "Ryn", in.Body.Name; want != got {
			t.Fatalf("expected name in request body '%s', instead got: '%s'", want, got)
		}

		return OKResult[Response, Err](
			Response{
				Name: "Corynth",
			},
			http.StatusOK,
		)
	})

	// Test inputs to our ServeHTTP func
	rec := httptest.NewRecorder()
	body := Input{
		Name: "Ryn",
	}
	req := httptest.NewRequest(http.MethodPost, "http://example.com", encode(body))

	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	// Should be a not found code and response
	if want, got := http.StatusOK, resp.StatusCode; want != got {
		t.Fatalf("expected status code %d, instead got: %d", want, resp.StatusCode)
	}
	response := decode[Response](resp.Body)
	if want, got := "Corynth", response.Name; want != got {
		t.Fatalf("expected '%s', instead got: '%s'", want, got)
	}
}
