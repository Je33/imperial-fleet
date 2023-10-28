package main

import (
	"log"
	"github.com/Je33/imperial_fleet/internal/transport/rest"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := rest.RunRest()
	if err != nil {
		return err
	}
	return nil
}