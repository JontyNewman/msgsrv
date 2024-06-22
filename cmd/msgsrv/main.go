package main

import (
	"log"
	"os"

	"internal/repo"

	"github.com/jontynewman/msgsrv"
)

func main() {

	args := os.Args[1:]
	max := 1
	len := len(args)

	if len > max {
		log.Fatalf("expected no more than %d command-line argument but got %d", max, len)
	}

	addr := ""

	if len > 0 {
		addr = args[0]
	}

	err := msgsrv.Run(addr, new(repo.RuntimeMessageRepository))

	log.Fatal(err)
}
