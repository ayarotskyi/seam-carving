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

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
