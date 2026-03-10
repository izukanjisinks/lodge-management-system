package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]interface{}{
		"error": map[string]string{
			"message": message,
		},
	})
}

func DecodeJson(r *http.Request, v interface{}) error {
	// Read the body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("ERROR: Failed to read request body: %v\n", err)
		return err
	}
	defer r.Body.Close()

	// Log the raw request body for debugging
	fmt.Printf("DEBUG: Raw request body: %s\n", string(bodyBytes))

	// Decode from the bytes
	if err := json.Unmarshal(bodyBytes, v); err != nil {
		fmt.Printf("ERROR: Failed to unmarshal JSON: %v\n", err)
		return err
	}

	// Restore the body so it can be read again if needed
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return nil
}
