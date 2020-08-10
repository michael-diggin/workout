package database

import (
	"database/sql"
	"log"

	"github.com/michael-diggin/workout"
)

// EnsureTableExists executes table creation query
func EnsureTableExists(db *sql.DB) {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS events
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user TEXT,
	sport TEXT,
	title TEXT,
	duration INTEGER
)`

// DBService implements the workout.EventService interface
type DBService struct {
	db *sql.DB
}

// NewDBService returns DBService
func NewDBService(db *sql.DB) *DBService {
	return &DBService{
		db: db,
	}
}

// Event returns the event by given ID
func (s *DBService) Event(id int) (*workout.Event, error) {
	var e workout.Event
	query := s.db.QueryRow(`SELECT * FROM events WHERE id=$1`, id)
	if err := query.Scan(&e.ID, &e.User, &e.Sport, &e.Title, &e.Duration); err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateEvent updates a given event
func (s *DBService) UpdateEvent(e *workout.Event) error {
	_, err :=
		s.db.Exec(`UPDATE events SET user=$1, sport=$2, title=$3, duration=$4 WHERE id=$5`,
			e.User, e.Sport, e.Title, e.Duration, e.ID)

	return err
}

// DeleteEvent removes from DB
func (s *DBService) DeleteEvent(id int) error {
	_, err := s.db.Exec(`DELETE FROM events WHERE id=$1`, id)
	return err
}

// CreateEvent makes a new event and adds to the DB
func (s *DBService) CreateEvent(e *workout.Event) (int, error) {
	sqlStatement := `INSERT INTO events (user, sport, title, duration)
	VALUES($1, $2, $3, $4)`
	res, err := s.db.Exec(sqlStatement, e.User, e.Sport, e.Title, e.Duration)
	id, _ := res.LastInsertId()
	e.ID = int(id)
	return int(id), err
}

// Events returns all events
func (s *DBService) Events() ([]workout.Event, error) {
	rows, err := s.db.Query(
		`SELECT id, user, sport, title, duration FROM events`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []workout.Event{}

	for rows.Next() {
		var e workout.Event
		if err := rows.Scan(&e.ID, &e.User, &e.Sport, &e.Title, &e.Duration); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}
