package mplex

import (
	"encoding/json"
	"net/http"
)

// Package level handlers
//
// While I like to avoid package level state/variables, this seemed the easiest
// way to provide good defaults to the consumer and while allowing them to configure
// routine things like error handling.
//
// If you find that one or a few of your routes does not benefit from this pattern,
// you'll be better off handling the exceptions manually in a vanilla http handler.
// That'll be more flexible than trying to make one function handle case.

var (
	// Encoder is used to output the result of the consumer-supplied handler.
	// It doesn't return an error since the function will handle errors itself.
	Encoder = JSONEncoder

	// Decoder is what the handler functions below use to decode requests with bodies.
	//
	// It returns if the operation was successful. It doesn't return an error since
	// the function will handle errors itself.
	Decoder = JSONDecoder

	// JSONEncoder writes the value out as JSON. If unable, it will write the error and
	// a 500 error code
	JSONEncoder = func(w http.ResponseWriter, out any, status int) bool {
		// We wrote JSON, so tell the response that it's coming back as JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		if err := json.NewEncoder(w).Encode(out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return false
		}

		return true
	}

	// JSONDecoder unmarshals the body into the given destination
	JSONDecoder = func(w http.ResponseWriter, r *http.Request, dest any) bool {
		if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return false
		}

		return true
	}
)

type Nillable interface {
}

// Result holds the necessary fields that will be output for a response
type Result[T, E any] struct {
	StatusCode int // If not set, this will be a 200: http.StatusOK

	value T
	err   *E
}

// OKResult constructs a Result from the value and status code, no error
// value added
func OKResult[T, E any](v T, statusCode int) Result[T, E] {
	return Result[T, E]{
		StatusCode: statusCode,
		value:      v,
	}
}

// ErrResult constructs a Result with the given error set
func ErrResult[T, E any](err E, statusCode int) Result[T, E] {
	return Result[T, E]{
		StatusCode: statusCode,
		err:        &err,
	}
}

// Request holds the request details for requests without a body, e.g. GET
//
// Instead of using a plain http.Request, this gives us the ability to
// add fields later and evolve this type without breaking the package's API.
type Request struct {
	Request *http.Request
}

// Handler is the type for a function that gets a request without a body
type Handler[Out, E any] func(i Request) Result[Out, E]

// ServeHTTP implements the http.Handler interface
func (h Handler[Out, E]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h(Request{
		Request: r,
	})

	var outVal any = res.value
	if res.err != nil {
		outVal = *res.err
	}

	// If there's a StatusCode, use that as the header
	status := http.StatusOK
	if res.StatusCode > 0 {
		status = res.StatusCode
	}

	// Write the value back out
	if ok := Encoder(w, outVal, status); !ok {
		return
	}
}

// BodyRequest is the decoded request with the associated body
type BodyRequest[T any] struct {
	Request *http.Request
	Body    T
}

// BodyHandler is a type that represents a handler that will recieve a body
// and output a body in the response.
type BodyHandler[In, Out, E any] func(i BodyRequest[In]) Result[Out, E]

// ServeHTTP implements the http.Handler interface
func (h BodyHandler[In, Out, E]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the body into the In type to pass into the function
	var in In
	if ok := Decoder(w, r, &in); !ok {
		// Decoder will write the status and error out, we just need to exit here
		return
	}

	res := h(BodyRequest[In]{
		Request: r,
		Body:    in,
	})

	var outVal any = res.value
	if res.err != nil {
		outVal = *res.err
	}

	// If there's a StatusCode, use that as the header
	status := http.StatusOK
	if res.StatusCode > 0 {
		status = res.StatusCode
	}

	// Write the value back out
	if ok := Encoder(w, outVal, status); !ok {
		return
	}
}
