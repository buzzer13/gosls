# GoSLS

[![latest release](https://img.shields.io/github/v/release/buzzer13/gosls)](https://github.com/buzzer13/gosls/releases)

This project provides structure definitions and helper methods for writing serverless functions in the Go programming language.

## Supported services

- [DigitalOcean Functions](https://docs.digitalocean.com/products/functions/)
    - Supported runtime versions: `go:1.20`

# Examples

```go
package main

import (
    "github.com/buzzer13/gosls/do"
    "net/http"
)

func Main(evm do.FuncEventMap) (do.FuncResponseMap, error) {
    // Parse incoming event map
    evt, err := evm.Event()

    if err != nil {
        return (&do.FuncResponse{
            Body:       "event parse error - " + err.Error(),
            StatusCode: http.StatusInternalServerError,
        }).Map(), err
    }

    // Generate http.Request from the Event object
    req, err := evt.Request()

    if err != nil {
        return (&do.FuncResponse{
            Body:       "request parse error - " + err.Error(),
            StatusCode: http.StatusBadRequest,
        }).Map(), err
    }

    // Create response writer structure
    res := do.FuncResponseWriter{}

    // Dispatch HTTP request to your app and write result to the response.
    // With Echo framework this may be:
    // echoInstance.ServeHTTP(&res, req)

    // Generate map from the response object and return it from the function
    return res.GetFuncResponse().Map(), nil
}
```

# Disclaimer

GoSLS is an unofficial library and is not affiliated with any of the companies whose services are supported by this library. All trademarks and registered trademarks are the property of their respective owners.
