package utils

import "log"

func ErrorPanic(err error) {
	if err != nil {
		log.Fatalf("System failed: %v", err)
	}
}
