package main

import (
	"fmt"
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
		Scope:   os.Getenv("UNIEX_SCOPE"),   // default: client [client, infra]
	}

	// perform Backup of all Appliances xml configuration
	out, err := c.Export()
	if err != nil {
		fmt.Printf(_app+"[ERROR][EXIT] %s\n", err)
		os.Exit(1)
	}
	fmt.Printf(out)
}
