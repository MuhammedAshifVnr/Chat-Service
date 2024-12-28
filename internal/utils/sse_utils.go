package utils

import (
	"fmt"
	"io"
)

// WriteSSE sends a single Server-Sent Event to the client.
func WriteSSE(w io.Writer, eventType, data string) error {
	// Format the SSE event data.
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, data)

	// Write the message to the writer (e.g., http.ResponseWriter for SSE).
	_, err := w.Write([]byte(message))
	return err
}

// WriteSSEJSON sends a JSON payload as a Server-Sent Event.
func WriteSSEJSON(w io.Writer, eventType string, jsonData []byte) error {
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, string(jsonData))
	_, err := w.Write([]byte(message))
	return err
}