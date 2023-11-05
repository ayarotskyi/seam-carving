package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	L "seam-carving/handlers"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/hello", L.Hello_handler)

	err := http.ListenAndServe(":4000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
