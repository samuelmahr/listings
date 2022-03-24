package main

import (
	"github.com/samuelmahr/listings/internal/application"
	"github.com/samuelmahr/listings/internal/configuration"
	"log"
)

func main() {
	c, err := configuration.Configure()
	if err != nil {
		log.Fatal(err)
	}

	app := application.NewAPIApplication(c)
	app.Run()
}
