package helpers

import (
	"log"
)

func Must(err error) {
	if err != nil {
		// log.Fatalf("Fatal error in Must: %v\n", err)
		log.Printf("Error in Must: %v\n", err)
	}
}
