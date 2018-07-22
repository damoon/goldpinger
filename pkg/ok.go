package goldpinger

import (
	"log"
	"net/http"
)

// OK confirms a http connection was created
func OK(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Printf("failed to send response: %v", err)
	}
}
