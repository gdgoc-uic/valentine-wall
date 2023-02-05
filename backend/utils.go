package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func passivePrintError(err error) {
	log.Println(err.Error())
}

func getTermsAndConditions() ([]byte, error) {
	return os.ReadFile(filepath.Join(dataDirPath, "terms-and-conditions.md"))
}

// sse-related
func encodeDataSSE(rw http.ResponseWriter, data interface{}) {
	writer := &bytes.Buffer{}
	encoder := json.NewEncoder(writer)
	if err := encoder.Encode(data); err != nil {
		log.Println(err)
		fmt.Fprintf(writer, "null")
	}
	writeResponseDataSSE(rw, writer)
}

func writeResponseDataSSE(rw http.ResponseWriter, buf *bytes.Buffer) {
	fmt.Fprint(rw, "data: ")
	buf.WriteTo(rw)
	fmt.Fprint(rw, "\n\n")
	if f, ok := rw.(http.Flusher); ok {
		f.Flush()
	}
}

func checkProfanity(content string) *ResponseError {
	if profanityDetector.IsProfane(content) {
		return &ResponseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Your submission contains inappropriate content.",
		}
	}

	return nil
}
