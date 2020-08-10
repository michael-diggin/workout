package main

import "errors"

// Event struct to hold one exercise data
type Event struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Sport    string `json:"sport"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

// AllEvents is the Mock database
type AllEvents struct {
	Events []Event
}

//GetAllEvents returns all the events in the database
func GetAllEvents(db *AllEvents) ([]Event, error) {
	return nil, errors.New("Not Implemented")
}
