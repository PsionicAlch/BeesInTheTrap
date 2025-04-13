package main

import (
	"log"
	"os"

	"github.com/PsionicAlch/BeesInTheTrap/internal/game"
)

func main() {
	communication := game.StartupServer()
	client := createClient(communication, os.Stdin, os.Stdout, func(err error) {
		log.Fatalln(err)
	})

	client.run()
}
