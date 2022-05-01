package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jdholdren/upfront"
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

// APIError is a structured output of an error that occurred
type APIError struct {
	Err error
}

// MarshalJSON is implemented so our APIError looks good returning to the client
func (e APIError) MarshalJSON() ([]byte, error) {
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

func HandleGetValue(db *MutexDB[Record]) upfront.Handler[Record, APIError] {
	return func(in upfront.Request) upfront.Result[Record, APIError] {
		vars := mux.Vars(in.Request)
		key := vars["key"]

		val, ok := db.Get(key)
		if !ok {
			// If we couldn't find it, return a 404
			return upfront.ErrResult[Record, APIError](
				APIError{
					Err: fmt.Errorf("key not found in the db: %s", key),
				},
				http.StatusNotFound,
			)
		}

		return upfront.OKResult[Record, APIError](
			val,
			http.StatusOK,
		)
	}
}

func HandleSetValue(db *MutexDB[Record]) upfront.BodyHandler[Record, Record, APIError] {
	return func(in upfront.BodyRequest[Record]) upfront.Result[Record, APIError] {
		vars := mux.Vars(in.Request)
		key := vars["key"]

		db.Set(key, in.Body)

		return upfront.OKResult[Record, APIError](
			in.Body,
			http.StatusOK,
		)
	}
}
