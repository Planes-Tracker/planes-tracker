package main

import (
	"fmt"
	"log"

	"github.com/LockBlock-dev/planes-tracker/app"
)

func main()  {
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize: %w", err))
	}
	defer app.Stop()

	app.Start()
}
