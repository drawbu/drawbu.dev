package app

import "fmt"

type App struct {
}

func (app *App) Run() {
	fmt.Printf("hello, world from app\n")
}
