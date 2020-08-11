package routing

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/michael-diggin/workout"
)

//App struct to hold router and Event service
// App handles all the routing
type App struct {
	Router  *mux.Router
	Service workout.EventService
}

// SetUp will initialize the application
func (a *App) SetUp(service workout.EventService) {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.Service = service
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/events", a.GetEvents).Methods("GET")
	a.Router.HandleFunc("/event", a.CreateEvent).Methods("POST")
	a.Router.HandleFunc("/event/{id}", a.GetEvent).Methods("GET")
	a.Router.HandleFunc("/event/{id}", a.UpdateEvent).Methods("PUT")
	a.Router.HandleFunc("/event/{id}", a.DeleteEvent).Methods("DELETE")
}

//Run will start the application
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

//ServeHTTP will serve and route a request
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Router.ServeHTTP(w, r)
}
