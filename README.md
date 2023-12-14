# Http Builder

This is a Go project that provides a builder for `http.Request`. It allows you to easily construct HTTP requests with
various methods, headers, body types, and query parameters.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing
purposes.

### Prerequisites

You need to have Go installed on your machine. You can download it from the official website:

```
https://golang.org/dl/
```

### Installing

```
go get -u https://github.com/eatmoreapple/httpx
```

## Usage

Import the `httpx` package in your Go file:

```go
import "github.com/eatmoreapple/httpx"
```

Create a new `RequestBuilder`:

```go
builder := httpx.New("http://example.com")
```

Set up the request:

```go
builder.Method("POST").SetHeader("Accept", "application/json")
```

Build and send the request:

```go
response, err := builder.Do()
```

## Running the tests

Explain how to run the automated tests for this system.