package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michael-diggin/workout"
)

func setUpDB(t *testing.T) (DBService, func(DBService)) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Could not open databse: %s\n", err)
	}
	EnsureTableExists(db)
	dbService := NewDBService(db)
	return *dbService, func(service DBService) { service.db.Close() }
}

func addEvent(t *testing.T, service *DBService) {
	_, err := service.db.Exec(`INSERT INTO events (user, sport, title, duration) VALUES($1, $2, $3, $4)`,
		"test event", "test run", "test title", 10)
	if err != nil {
		t.Errorf("Could not add event to DB: %s\n", err)
	}
	return
}

func TestEnsureTableExists(t *testing.T) {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("Could not open databse: %s\n", err)
	}
	EnsureTableExists(db)
	_, err = db.Exec(`SELECT * FROM events`)
	if err != nil {
		t.Fatalf("Could not access 'events' table: %v", err)
	}
}

func TestEmptyTable(t *testing.T) {
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)

	allEvents, err := dbService.Events()
	if err != nil {
		t.Errorf("Error getting all events: %s\n", err)
	}
	if len(allEvents) > 0 {
		t.Fatalf("Expected an empty array. Got %v", allEvents)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)

	_, err := dbService.Event(1)
	if err != sql.ErrNoRows {
		t.Fatalf("Expected sql No rows error, got %s\n", err)
	}
}

func TestCreateEvent(t *testing.T) {
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)

	e := workout.Event{User: "test", Sport: "test run", Title: "test title", Duration: 1}
	id, err := dbService.CreateEvent(&e)
	if err != nil {
		t.Fatalf("Error creating event: %s\n", err)
	}

	if id != 1 {
		t.Fatalf("Expected event ID to be '1'. Got '%v'", id)
	}
}

func TestGetEvent(t *testing.T) {
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)
	addEvent(t, &dbService)

	event, err := dbService.Event(1)
	if err != nil {
		t.Errorf("Error getting event: %s\n", err)
	}
	if event.User != "test event" {
		t.Fatalf("Expected 'test event', got %v", event.User)
	}
}

func TestUpdateEvent(t *testing.T) {
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)
	addEvent(t, &dbService)

	e := workout.Event{ID: 1, User: "test-updated", Sport: "test run", Title: "test title", Duration: 11}

	err := dbService.UpdateEvent(&e)
	if err != nil {
		t.Errorf("Error updating event: %s\n", err)
	}

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
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)
	addEvent(t, &dbService)

	err := dbService.DeleteEvent(1)
	if err != nil {
		t.Fatalf("Error deleting event: %s\n", err)
	}
}

func TestGetEvents(t *testing.T) {
	dbService, teardown := setUpDB(t)
	defer teardown(dbService)
	addEvent(t, &dbService)
	addEvent(t, &dbService) // add two events

	allEvents, err := dbService.Events()
	if err != nil {
		t.Fatalf("Could not get all events: %v", err)
	}
	if len(allEvents) != 2 {
		t.Errorf("Expected 2 events, got %v", len(allEvents))
	}
	e := allEvents[0]
	if e.ID != 1 {
		t.Fatalf("Expected event with ID 1, got %v", e.ID)
	}
}
