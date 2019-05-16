# Golang Problem Details

[![Build Status](https://travis-ci.com/mvmaasakkers/go-problemdetails.svg?branch=master)](https://travis-ci.com/mvmaasakkers/go-problemdetails) 
[![MIT license](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/mvmaasakkers/go-problemdetails?status.svg)](https://godoc.org/github.com/mvmaasakkers/go-problemdetails)
[![Go Report Card](https://goreportcard.com/badge/github.com/mvmaasakkers/go-problemdetails)](https://goreportcard.com/report/github.com/mvmaasakkers/go-problemdetails)
[![CodeFactor](https://www.codefactor.io/repository/github/mvmaasakkers/go-problemdetails/badge)](https://www.codefactor.io/repository/github/mvmaasakkers/go-problemdetails)
[![Coverage Status](https://coveralls.io/repos/github/mvmaasakkers/go-problemdetails/badge.svg?branch=master)](https://coveralls.io/github/mvmaasakkers/go-problemdetails?branch=master)

Problem details implementation (https://tools.ietf.org/html/rfc7807) package for go.

`go get github.com/mvmaasakkers/go-problemdetails`

## How to use

The `ProblemDetails` struct can be used as `error` because it implements the `error` interface. The `ProblemType`
interface can be used to create predefined `ProblemDetails` with extensions and also implements the `error` interface.

The struct is setup to be used by the [json](https://golang.org/pkg/encoding/json/) and 
[xml](https://golang.org/pkg/encoding/xml/) marshaler from the stdlib and will marshal into `application/problem+json`
or `application/problem+xml` compliant data as defined in the [RFC 7807](https://tools.ietf.org/html/rfc7807).

To generate a `ProblemDetails` based on just a HTTP Status Code you can create one using `NewHTTP(statusCode int)`:

```go
problemDetails := problemdetails.NewHTTP(http.StatusNotFound)
``` 

This will generate a `ProblemDetails` struct that marshals as follows:

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404
}
```

```xml
<problem xmlns="urn:ietf:rfc:7807">
    <type>about:blank</type>
    <title>Not Found</title>
    <status>404</status>
</problem>
```

or use the more verbose `New(statusCode int, problemType, title, detail, instance string)`:

```go
problemDetails := problemdetails.New(http.StatusNotFound, "https://example.net/problem/object_not_found", "Object not found", "Object with id 1234 was not found, another id should be given.", "https://api.example.net/objects/1234")
``` 

This will generate a `ProblemDetails` struct that marshals as follows:

```json
{
  "type": "https://example.net/problem/object_not_found",
  "title": "Object not found",
  "status": 404,
  "detail": "Object with id 1234 was not found, another id should be given.",
  "instance": "https://api.example.net/objects/1234"
}
```

```xml
<problem xmlns="urn:ietf:rfc:7807">
    <type>https://example.net/problem/object_not_found</type>
    <title>Object not found</title>
    <status>404</status>
    <detail>Object with id 1234 was not found, another id should be given.</detail>
    <instance>https://api.example.net/objects/1234</instance>
</problem>
```

## Http helpers

For ease of use there are two output handlers available. `ProblemDetails.ServeJSON` for JSON and `ProblemDetails.ServeXML` for XML.

A shorthand for generating a 404 statuscode in Problem Details JSON to the ResponseWriter you can:

```go
problemdetails.NewHTTP(http.StatusNotFound).ServeJSON(w, r)
```

 