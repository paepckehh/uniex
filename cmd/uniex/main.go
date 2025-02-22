package main

import (
	"fmt"
	"log"
	"os"

	"paepcke.de/uniex"
)

const (
	_app = "[UNIEX]"
)

func main() {

	// setup from env
	c := &uniex.Config{
		MongoDB: os.Getenv("UNIEX_MONGODB"), // default: mongodb://localhost:27117
		Format:  os.Getenv("UNIEX_FORMAT"),  // default: csv [csv, json]
		Scope:   os.Getenv("UNIEX_SCOPE"),   // default: client [client] TODO: infra
	}

	// perform Backup of all Appliances xml configuration
	out, err := c.Export()
	if err != nil {
		log.Fatalf(_app+"[ERROR][EXIT] %v\n", err)
	}
	fmt.Print(string(out))
}
