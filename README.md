# Backoff [![GoDoc][godoc image]][godoc]

This is a Go port of the exponential backoff algorithm modified from [GRPC-go lib][grpc-go-backoff].

[Exponential backoff][exponential backoff wiki]
is an algorithm that uses feedback to multiplicatively decrease the rate of some process,
in order to gradually find an acceptable rate.
The retries exponentially increase and stop increasing when a certain threshold is met.

## Usage

`go get -u github.com/lixin9311/backoff`.

Use https://pkg.go.dev/github.com/lixin9311/backoff to view the documentation.

## Contributing

- I would like to keep this library as small as possible.
- Please don't send a PR without opening an issue and discussing it first.
- If proposed change is not a common use case, I will probably not accept it.

[godoc]: https://pkg.go.dev/github.com/lixin9311/backoff
[godoc image]: https://pkg.go.dev/github.com/lixin9311/backoff?status.png
[grpc-go-backoff]: https://pkg.go.dev/google.golang.org/grpc/backoff
[exponential backoff wiki]: http://en.wikipedia.org/wiki/Exponential_backoff
[advanced example]: https://pkg.go.dev/github.com/cenkalti/backoff/v4?tab=doc#pkg-examples
