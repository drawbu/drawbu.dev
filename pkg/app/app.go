package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type App struct {
	port int16
}

func New(port int16) *App {
	return &App{port}
}

func (app *App) Run() {
	component := hello("John")

	http.Handle("/", templ.Handler(component))

	fmt.Printf("Listening on :%d\n", app.port)
    err := http.ListenAndServe(":" + strconv.Itoa(int(app.port)), nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
}
