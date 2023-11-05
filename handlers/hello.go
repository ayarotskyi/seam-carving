package handlers

import (
	"encoding/json"
	"net/http"
)

func Hello_handler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello world")
}
