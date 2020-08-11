package routing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/michael-diggin/workout"
	"github.com/michael-diggin/workout/mock"
)

//TODO test table to get all the different errors

var testEvent = workout.Event{
	ID:       1,
	User:     "test user",
	Sport:    "test sport",
	Title:    "test title",
	Duration: 10,
}

func TestGetEvent(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	ser.EventFn = func(id int) (*workout.Event, error) {
		if id != 1 {
			return nil, sql.ErrNoRows
		}
		return &testEvent, nil
	}

	req, _ := http.NewRequest("GET", "/event/1", nil)
	rec := httptest.NewRecorder()
	a.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}
	var event map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &event)
	if event["user"] != "test user" {
		t.Errorf("Expected user = 'test user', got %v", event["user"])
	}
	if event["id"] != 1.0 {
		t.Errorf("Expected id = '1', got %v", event["id"])
	}
	if !ser.EventInvoked {
		t.Errorf("Expected Event() to be invoked")
	}
}

func TestGetAllEvents(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	ser.EventsFn = func() ([]*workout.Event, error) {
		allEvents := []*workout.Event{
			&testEvent,
		}
		return allEvents, nil
	}

	req, _ := http.NewRequest("GET", "/events", nil)
	rec := httptest.NewRecorder()
	a.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}

	if !ser.EventsInvoked {
		t.Errorf("Expected Events() to be invoked")
	}
}

func TestCreateEvent(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	ser.CreateFn = func(e *workout.Event) (int, error) {
		return 10, nil
	}

	var jsonStr = []byte(`{"user":"test", "sport": "test run", "title": "test title", "duration": 1}`)
	req, _ := http.NewRequest("POST", "/event", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	a.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected 201 response code. Got %d\n", rec.Code)
	}
	var event map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &event)
	if event["id"] != 10.0 {
		t.Errorf("Expected 'id=10' but got %v", event["id"])
	}
	if !ser.CreateInvoked {
		t.Errorf("Expected CreateEvent() to be invoked")
	}
}

func TestUpdateEvent(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	ser.UpdateFn = func(e *workout.Event) error {
		return nil
	}

	var jsonStr = []byte(`{"user":"testupdate", "sport": "test run", "title": "test title", "duration": 1}`)
	req, _ := http.NewRequest("PUT", "/event/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	a.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}
	var event map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &event)
	if event["id"] != 1.0 {
		t.Errorf("Expected 'id=10' but got %v", event["id"])
	}
	if event["user"] != "testupdate" {
		t.Errorf("Expected 'user=testupdate' but got %v", event["id"])
	}
	if !ser.UpdateInvoked {
		t.Errorf("Expected CreateEvent() to be invoked")
	}
}

func TestDeleteEvent(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	ser.DeleteFn = func(id int) error {
		return nil
	}
	req, _ := http.NewRequest("DELETE", "/event/1", nil)
	rec := httptest.NewRecorder()
	a.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200 response code. Got %d\n", rec.Code)
	}
	if !ser.DeleteInvoked {
		t.Errorf("Expected CreateEvent() to be invoked")
	}
}
