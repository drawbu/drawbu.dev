package app

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

type App struct {
}

func (app *App) Run() {
	component := hello("John")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
