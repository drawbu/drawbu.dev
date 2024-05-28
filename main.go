package main

import (
	"server/pkg/app"
)

func main() {
	server := app.App{Port: 8080}
	server.Run()
}
