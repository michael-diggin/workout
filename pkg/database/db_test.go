package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michael-diggin/workout"
)

var dbService DBService

func TestMain(m *testing.M) {
	setUpDB()
	defer dbService.db.Close()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func setUpDB() {
	db, err := Open(":memory:")
	if err != nil {
		log.Fatalf("Could not open databse: %s\n", err)
	}
	EnsureTableExists(db)
	dbService.db = db
}

func clearTable() {
	dbService.db.Exec("DELETE FROM events")
	dbService.db.Exec("DELETE FROM sqlite_sequence WHERE name='events';")
}

func addEvent(t *testing.T) {
	_, err := dbService.db.Exec(`INSERT INTO events (user, sport, title, duration) VALUES($1, $2, $3, $4)`,
		"test event", "test run", "test title", 10)
	if err != nil {
		t.Errorf("Could not add event to DB: %s\n", err)
	}
	return
}

// TODO test table creation/ensure table exists function

func TestEmptyTable(t *testing.T) {
	clearTable()

	allEvents, err := dbService.Events()
	if err != nil {
		t.Errorf("Error getting all events: %s\n", err)
	}
	if len(allEvents) > 0 {
		t.Errorf("Expected an empty array. Got %v", allEvents)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	clearTable()

	_, err := dbService.Event(1)
	if err != sql.ErrNoRows {
		t.Fatalf("Expected sql No rows error, got %s\n", err)
	}
}

func TestCreateEvent(t *testing.T) {
	clearTable()

	e := workout.Event{User: "test", Sport: "test run", Title: "test title", Duration: 1}
	id, err := dbService.CreateEvent(&e)
	if err != nil {
		t.Errorf("Error creating event: %s\n", err)
	}

	if id != 1 {
		t.Errorf("Expected event ID to be '1'. Got '%v'", id)
	}
}

func TestGetEvent(t *testing.T) {
	clearTable()
	addEvent(t)

	event, err := dbService.Event(1)
	if err != nil {
		t.Errorf("Error getting event: %s\n", err)
	}
	if event.User != "test event" {
		t.Errorf("Expected 'test event', got %v", event.User)
	}
}

func TestUpdateEvent(t *testing.T) {
	clearTable()
	addEvent(t)

	e := workout.Event{ID: 1, User: "test-updated", Sport: "test run", Title: "test title", Duration: 11}

	err := dbService.UpdateEvent(&e)
	if err != nil {
		t.Errorf("Error updating event: %s\n", err)
	}
	// get the same event and check is has been updates
	newEvent, err := dbService.Event(1)
	if err != nil {
		t.Errorf("Error getting updated event: %s\n", err)
	}

	if newEvent.User != "test-updated" {
		t.Errorf("Expected updated name (%v). Got %v", "test-updated", newEvent.User)
	}

	if newEvent.Duration != 11 {
		t.Errorf("Expected updated duration '11'. Got '%v'", newEvent.Duration)
	}
}

func TestDeleteEvent(t *testing.T) {
	clearTable()
	addEvent(t)

	err := dbService.DeleteEvent(1)
	if err != nil {
		t.Errorf("Error deleting event: %s\n", err)
	}
}
