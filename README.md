# go-common

`go-common` is a collection of reusable functions and utilities for web services.

**Context**
```go
package main

import (
    "github.com/h-okay/go-common"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    session := common.ContextKey("session")
    common.ContextSet(r, session, "foobar")

    value := common.ContextGet[string](r, session)
    w.Write([]byte(value))
}
```

**Tokens**
```go
token, err := common.ExtractBearerToken(r)
if err != nil {
    http.Error(w, err.Error(), http.StatusUnauthorized)
    return
}
```

**Logging**
```go
logger := common.NewLogger(os.Stdout, common.LevelInfo)

logger.PrintInfo("This is an info message", nil)
logger.PrintError(errors.New("An error occurred"), nil)
```

**JSON**
```go
type Response struct {
    Message string `json:"message"`
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
    var data Response
    if err := common.ReadJSON(w, r, &data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    response := Envelope{"message": "Hello, " + data.Message}
    common.WriteJSON(w, http.StatusOK, response, nil)
}
```

## Installation

To install `go-common`, run:

```bash
go get github.com/h-okay/go-common
