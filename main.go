package main

import (
	"server/pkg/app"
)

func main() {
	server := app.New(8080)
	server.Run()
}
