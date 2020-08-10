package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.SetUp(":memory:")
	ensureTableExists(a.DB)
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func clearTable() {
	a.DB.Exec("DELETE FROM events")
	a.DB.Exec("DELETE FROM sqlite_sequence WHERE name='events';")
}

func addEvent(t *testing.T) {
	_, err := a.DB.Exec(`INSERT INTO events (user, sport, title, duration) VALUES($1, $2, $3, $4)`,
		"test event", "test run", "test title", 10)
	if err != nil {
		t.Errorf("Could not add event to DB: %s\n", err)
	}
	return
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/events", nil)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rr.Code)
	}

	if body := rr.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/event/11", nil)
	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404 response code. Got %d\n", rec.Code)
	}

	var m map[string]string
	json.Unmarshal(rec.Body.Bytes(), &m)
	if m["error"] != "Event not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Event not found'. Got '%s'", m["error"])
	}
}

func TestCreateEvent(t *testing.T) {

	clearTable()

	var jsonStr = []byte(`{"user":"test", "sport": "test run", "title": "test title", "duration": 1}`)
	req, _ := http.NewRequest("POST", "/event", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected 201 response code. Got %d\n", rec.Code)
	}

	var m map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["user"] != "test" {
		t.Errorf("Expected user name to be 'test'. Got '%v'", m["sser"])
	}
	if m["sport"] != "test run" {
		t.Errorf("Expected sport to be 'test run'. Got '%v'", m["sport"])
	}
	if m["title"] != "test title" {
		t.Errorf("Expected title to be 'test title'. Got '%v'", m["title"])
	}
	if m["duration"] != 1.0 {
		t.Errorf("Expected duration to be '1'. Got '%v'", m["duration"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected event ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetEvent(t *testing.T) {
	clearTable()
	addEvent(t)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}
	var event map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &event)
	if event["user"] != "test event" {
		t.Errorf("Expected 'test', got %v", event["user"])
	}
}

func TestUpdateEvent(t *testing.T) {
	clearTable()
	addEvent(t)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)

	var originalEvent map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &originalEvent)

	var jsonStr = []byte(`{"name":"test1-updated name", "sport": "test run", "title": "test title", "duration": 11}`)
	req, _ = http.NewRequest("PUT", "/event/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	rec = httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}

	var m map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &m)

	if m["id"] != originalEvent["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalEvent["id"], m["id"])
	}

	if m["user"] == originalEvent["user"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalEvent["user"], m["user"], m["user"])
	}

	if m["duration"] == originalEvent["duration"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalEvent["duration"], m["duration"], m["duration"])
	}
}

func TestDeleteEvent(t *testing.T) {
	clearTable()
	addEvent(t)

	req, _ := http.NewRequest("GET", "/event/1", nil)
	rec := httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}

	req, _ = http.NewRequest("DELETE", "/event/1", nil)
	rec = httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}

	req, _ = http.NewRequest("GET", "/event/1", nil)
	rec = httptest.NewRecorder()
	a.Router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404 response code. Got %d\n", rec.Code)
	}
}
