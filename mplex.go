package mplex

import (
	"encoding/json"
	"net/http"
)

// Err is a type constraint that requires that we pass a pointer
// to our type that will represent an error
type Err[T any] interface {
	*T
}

// Req is what comes into a handler function
type Req[T any] struct {
	Request *http.Request
	Body    T
}

// Result holds the necessary fields that will be output for a response
type Result[T, E any, ErrT Err[E]] struct {
	Value      T
	Err        ErrT
	StatusCode int // If not set, this will be a 200: http.StatusOK
}

type InOutHandler[In, Out, E any] func(i Req[In]) Result[Out, E, *E]

func (h InOutHandler[In, Out, E]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the body into the In type to pass into the function
	var in In
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		// TODO: Need to handle this error custom-ly according to the client
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := h(Req[In]{
		Request: r,
		Body:    in,
	})

	// If there's a StatusCode, use that as the header
	if res.StatusCode > 0 {
		w.WriteHeader(res.StatusCode)
	}

	var outVal any = res.Value
	if res.Err != nil {
		outVal = res.Err
	}

	w.Header().Set("Content-Type", "application/json")

	// Write the value back out
	if err := json.NewEncoder(w).Encode(outVal); err != nil {
		// TODO: Need to handle this error custom-ly according to the client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type OutHandler[Out, E any] func(i Req[struct{}]) Result[Out, E, *E]

func (h OutHandler[Out, E]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := h(Req[struct{}]{
		Request: r,
	})

	// If there's a StatusCode, use that as the header
	if res.StatusCode > 0 {
		w.WriteHeader(res.StatusCode)
	}

	var outVal any = res.Value
	if res.Err != nil {
		outVal = res.Err
	}

	w.Header().Set("Content-Type", "application/json")

	// Write the value back out
	if err := json.NewEncoder(w).Encode(outVal); err != nil {
		// TODO: Need to handle this error custom-ly according to the client
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
