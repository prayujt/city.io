package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Serving at localhost:8000...")
	http.ListenAndServe(":8000", nil)
}
