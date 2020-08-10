package main

import "github.com/gorilla/mux"

//App struct to hold router and database
type App struct {
	Router *mux.Router
	DB     AllEvents
}

// SetUp will initialize the application
func (a *App) SetUp() {
	a.Router = mux.NewRouter()
	event := Event{1, "Michael", "Cycling", "Quick Spin", 60}

	a.DB = AllEvents{[]Event{
		event,
	}}
}

//Run will start the application
func (a *App) Run(addr string) {}
