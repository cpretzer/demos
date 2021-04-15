package main

import (
	"log"
	"os"

	"github.com/google/go-github/v34/github"
)

func main() {
	github.NewClient(nil)
	log.Println("doing github things")
	os.Exit(0)
}
