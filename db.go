package main

import (
	"database/sql"
)

// Event struct to hold one exercise data
type Event struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Sport    string `json:"sport"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

// GetEvent returns the event by given ID
func (e *Event) GetEvent(db *sql.DB) error {
	query := db.QueryRow(`SELECT * FROM events WHERE id=$1`, e.ID)
	return query.Scan(&e.User, &e.Sport, &e.Title, &e.Duration)
}

//UpdateEvent updates a given event
func (e *Event) UpdateEvent(db *sql.DB) error {
	_, err :=
		db.Exec(`UPDATE events SET user=$1, sport=$2, title=$3, duration=$4 WHERE id=$5`,
			e.User, e.Sport, e.Title, e.Duration, e.ID)

	return err
}

//DeleteEvent removes from DB
func (e *Event) DeleteEvent(db *sql.DB) error {
	_, err := db.Exec(`DELETE FROM events WHERE id=$1`, e.ID)

	return err
}

//CreateEvent makes a new event and adds to the DB
func (e *Event) CreateEvent(db *sql.DB) error {
	sqlStatement := `INSERT INTO events (user, sport, title, duration)
	VALUES($1, $2, $3, $4)`
	res, err := db.Exec(sqlStatement, e.User, e.Sport, e.Title, e.Duration)
	id, _ := res.LastInsertId()
	e.ID = int(id)
	return err
}

//GetEvents returns all events
func GetEvents(db *sql.DB, start, count int) ([]Event, error) {
	rows, err := db.Query(
		"SELECT id, user, sport, title, duration FROM events LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []Event{}

	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.User, &e.Sport, &e.Title, &e.Duration); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}
