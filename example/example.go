package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jdholdren/mplex"
)

func main() {
	db := &MutexDB[Record]{
		Data: map[string]Record{},
	}

	db.Set("james", Record{
		Name: "james",
	})

	mux := mux.NewRouter()
	mux.Handle("/{key}", HandleSetValue(db)).Methods(http.MethodPut)
	mux.Handle("/{key}", HandleGetValue(db)).Methods(http.MethodGet)

	fmt.Println("Serving on port 4444...")
	if err := http.ListenAndServe(":4444", mux); err != nil {
		log.Fatalf("error service: %s", err)
	}
}

type Record struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// APIError is a structure output of an error that occurred
type APIError struct {
	Err error
}

// MarshalJSON is implemented so our APIError looks good returning to the client
func (e *APIError) MarshalJSON() ([]byte, error) {
	msg := ""
	if e.Err != nil {
		msg = e.Err.Error()
	}

	byts, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: msg,
	})
	if err != nil {
		return nil, err
	}

	return byts, nil
}

func HandleSetValue(db *MutexDB[Record]) mplex.InOutHandler[Record, Record, APIError] {
	return func(in mplex.Req[Record]) mplex.Result[Record, APIError, *APIError] {
		vars := mux.Vars(in.Request)
		key := vars["key"]

		db.Set(key, in.Body)

		return mplex.Result[Record, APIError, *APIError]{
			Value: in.Body,
		}
	}
}

func HandleGetValue(db *MutexDB[Record]) mplex.OutHandler[Record, APIError] {
	return func(in mplex.Req[struct{}]) mplex.Result[Record, APIError, *APIError] {
		vars := mux.Vars(in.Request)
		key := vars["key"]

		val, ok := db.Get(key)
		if !ok {
			// If we couldn't find it, return a 404
			return mplex.Result[Record, APIError, *APIError]{
				Err: &APIError{
					Err: fmt.Errorf("key not found in the db: %s", key),
				},
				StatusCode: http.StatusNotFound,
			}
		}

		return mplex.Result[Record, APIError, *APIError]{
			Value: val,
		}
	}
}
