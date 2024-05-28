package app

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"server/pkg/homepage"
)

type App struct {
	Port int16
}

func (app *App) Run() {
	// Routing
	http.HandleFunc("/", homepage.Handler)

	fmt.Printf("Listening on localhost:%d\n", app.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(int(app.Port)), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
