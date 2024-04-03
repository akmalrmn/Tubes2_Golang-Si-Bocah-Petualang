package handlers

import (
    "net/http"
    "fmt"
)

func SampleHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}
