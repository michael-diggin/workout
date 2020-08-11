package routing

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/michael-diggin/workout"
)

// GetEvent handles the Retrieve for /event/id
func (a *App) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var event *workout.Event
	event, err = a.Service.Event(id)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Event not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, event)
}

// GetEvents returns all events via /events endpoint
func (a *App) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, err := a.Service.Events()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, events)
}

// CreateEvent creates a new event and adds to the Database/Service
func (a *App) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var e *workout.Event
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	id, err := a.Service.CreateEvent(e)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//attach the new id to the event
	e.ID = id

	respondWithJSON(w, http.StatusCreated, e)
}

// UpdateEvent updates an event that already exists
func (a *App) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	var e *workout.Event
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	e.ID = id

	if err := a.Service.UpdateEvent(e); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, e)
}

// DeleteEvent removes the event from the DB/Service
func (a *App) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Event ID")
		return
	}

	if err := a.Service.DeleteEvent(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
