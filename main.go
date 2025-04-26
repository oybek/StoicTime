package main

import (
	"fmt"
	"log"

	"oybek.io/kerege/config"
)

func main() {
	theConfig, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error create config: %v", err)
	}

	fmt.Printf("%#v\n", theConfig)
}
