// Welcome mortal.
// This codebase was intentionally thrown together as a mess, with disregard to code quality.
// There are vulnerbilities *intentionally* left in for the students of room 235 to discover. We made this for an educational setting (to provide a small amount of utility and learning), hence what I stated before. 
// If you were thinking of using this as a reference for any sort of production code, don't. *This is only supposted to be run in a trusted, private enviroment.*
// You have been warned... have fun learning, child.

package main

import (
	"os"

	"github.com/gocopper/copper"
)

func main() {
	var app = copper.New()

	server, err := InitServer(app)
	if err != nil {
		app.Logger.Error("Failed to init server", err)
		os.Exit(1)
	}

	app.Start(server)
}
