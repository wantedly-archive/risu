package main

import (
	"log"

	"github.com/wantedly/risu/schema"
)

func printLog(build schema.Build, text string) {
	log.Printf("[%s] %s\n", build.ID.String(), text)
}
