package mock

import "github.com/michael-diggin/workout"

// TestService represents a mock implementation of myapp.UserService.
type TestService struct {
	EventFn      func(id int) (*workout.Event, error)
	EventInvoked bool

	EventsFn      func() ([]*workout.Event, error)
	EventsInvoked bool

	CreateFn      func(e *workout.Event) (int, error)
	CreateInvoked bool

	UpdateFn      func(e *workout.Event) error
	UpdateInvoked bool

	DeleteFn      func(id int) error
	DeleteInvoked bool
}

// Event invokes the mock implementation and marks the function as invoked.
func (s *TestService) Event(id int) (*workout.Event, error) {
	s.EventInvoked = true
	return s.EventFn(id)
}

// Events invokes the mock implementation and marks the function as invoked.
func (s *TestService) Events() ([]*workout.Event, error) {
	s.EventsInvoked = true
	return s.EventsFn()
}

// CreateEvent invokes the mock implementation and marks the function as invoked.
func (s *TestService) CreateEvent(e *workout.Event) (int, error) {
	s.CreateInvoked = true
	return s.CreateFn(e)
}

// UpdateEvent invokes the mock implementation and marks the function as invoked.
func (s *TestService) UpdateEvent(e *workout.Event) error {
	s.UpdateInvoked = true
	return s.UpdateFn(e)
}

// DeleteEvent invokes the mock implementation and marks the function as invoked.
func (s *TestService) DeleteEvent(id int) error {
	s.DeleteInvoked = true
	return s.DeleteFn(id)
}
