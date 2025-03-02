package main

import (
	"fmt"
	"log"
	"os"

	"paepcke.de/uniex"
)

const (
	_app = "UNIEX"
)

func main() {

	// setup from env
	c := &uniex.Config{
		MongoDB: os.Getenv(_app + "_MONGODB"), // default: mongodb://localhost:27117
		Format:  os.Getenv(_app + "_FORMAT"),  // default: csv [csv, json]
		Scope:   os.Getenv(_app + "_SCOPE"),   // default: client [client] TODO: infra
	}

	// perform Backup of all Appliances xml configuration
	out, err := c.Export()
	if err != nil {
		log.Fatalf("["+_app+"][ERROR][EXIT] %v\n", err)
	}
	fmt.Print(string(out))
}
