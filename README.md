# Upfront

Upfront is a very small library for adding type signatures to your HTTP handlers
using generics.

## Useage

For requests where you don't care about the request body:

```go
func handleGet() upfront.Handler[ResponseType, ErrorType] {
	return func(in upfront.Request) upfront.Result[ResponseType, ErrorType] {
		// `in` will contain your normall *http.Request

		// Use `upfront.OKResult` or `upfront.ErrResult` to return a status code
		// along with either the left or right type paramters
		return upfront.OKResult[ResponseType, ErrorType](
			val,
			http.StatusOK,
		)
	}
}
```

For requests where you want the request body, `upfront` will try to decode it
for you and pass it along:

```go
func handlePost() upfront.BodyHandler[Body, ResponseType, ErrorType] {
	return func(in upfront.BodyRequest[Record]) upfront.Result[Record, APIError] {
		// `in` not only has your `*http.Request`, but also has your `Body` type
		// in there as well

		// Do something with `in.Body`...

		// Maybe something has gone wrong
		return upfront.Err[Body, ResponseType, ErrorType](
			ErrorType{
				Err: err,
			},
			http.StatusInternalServerError,
		)
	}
}
```

## Example server

There's an example server and in-memory database to see how someone might use
the package.
Run `make example` to start the server at `:4444`.

## Contributing

This is a small library, so I don't imagine there will be too much to add to it,
but you're more than welcome to open a PR or open an issue!
If you're looking to get into the source code, you'll just need
[Go](https://go.dev/), and make sure it's at least version `1.18` to use
generics.
You can run `make test` at the root of the repo to easily run tests.
