package utils

import (
	"log"
	"strings"
)

func ErrorPanic(err error) {
	if err != nil {
		log.Fatalf("System failed: %v", err)
	}
}

func HandleConnectionClosed(err error, handler func()) {
	if err != nil && strings.Contains(err.Error(), "use of closed network connection") {
		log.Fatalf("Connection closed")
		handler()
	}
}
