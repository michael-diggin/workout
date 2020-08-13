package routing

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
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
		switch id {
		case 1:
			return &testEvent, nil
		case 2:
			return nil, errors.New("internal error")
		default:
			return nil, sql.ErrNoRows
		}
	}

	tt := []struct {
		name string
		id   string
		code int
	}{
		{"get normal event", "1", 200},
		{"get non existent event", "8", 404},
		{"bad request", "a", 400},
		{"internal error", "2", 500},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/event/"+tc.id, nil)
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("Expected %d response code. Got %d\n", tc.code, rec.Code)
			}
			if tc.code == http.StatusOK {
				var event map[string]interface{}
				json.Unmarshal(rec.Body.Bytes(), &event)
				if event["user"] != "test user" {
					t.Errorf("Expected user = 'test user', got %v", event["user"])
				}
				if !ser.EventInvoked {
					t.Errorf("Expected Event() to be invoked")
				}
			}
		})
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
		if e.ID != 1 {
			return errors.New("Non existing event")
		}
		return nil
	}

	tt := []struct {
		name string
		id   string
		code int
	}{
		{"update normal event", "1", 200},
		{"non existing event", "8", 500},
		{"bad request", "a", 400},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var jsonStr = []byte(`{"user":"testupdate", "sport": "test run", "title": "test title", "duration": 1}`)
			req, _ := http.NewRequest("PUT", "/event/"+tc.id, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("Expected %d response code. Got %d\n", tc.code, rec.Code)
			}
			if tc.code == 200 {
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
		})
	}
}

func TestBadPayload(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	tt := []struct {
		name   string
		method string
		url    string
	}{
		{"update", "PUT", "/event/1"},
		{"create", "POST", "/event"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var jsonStr = []byte(`{"User":10, "sport": , "title": "test title", "duration": 1}`)
			req, _ := http.NewRequest(tc.method, tc.url, bytes.NewBuffer(jsonStr))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("Expected 400 response code. Got %d\n", rec.Code)
			}
			var m map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &m)
			if m["error"] != "Invalid request payload" {
				t.Fatalf("expect bad payload err, got %s", m["error"])
			}
		})
	}

}

func TestDeleteEvent(t *testing.T) {
	var a App
	ser := &mock.TestService{}
	a.SetUp(ser)

	ser.DeleteFn = func(id int) error {
		if id != 1 {
			return errors.New("Does not exist")
		}
		return nil
	}

	tt := []struct {
		name string
		id   string
		code int
	}{
		{"delete normal event", "1", 200},
		{"delete non existent event", "8", 500},
		{"bad request", "a", 400},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			req, _ := http.NewRequest("DELETE", "/event/"+tc.id, nil)
			rec := httptest.NewRecorder()
			a.ServeHTTP(rec, req)

			if rec.Code != tc.code {
				t.Errorf("Expected %d response code. Got %d\n", tc.code, rec.Code)
			}
			if tc.id != "a" {
				if !ser.DeleteInvoked {
					t.Errorf("Expected CreateEvent() to be invoked")
				}
			}
		})
	}
}
