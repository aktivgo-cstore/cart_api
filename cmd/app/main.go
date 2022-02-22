package main

import (
	"cart_api/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalln("Unable to start app")
	}
}
