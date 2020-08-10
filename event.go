package workout

// Event struct to hold one exercise data
type Event struct {
	ID       int    `json:"id"`
	User     string `json:"user"`
	Sport    string `json:"sport"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
}

// EventService interface definition
type EventService interface {
	Event(id int) (*Event, error)
	Events() ([]*Event, error)
	CreateEvent(u *Event) (int, error)
	UpdateEvent(u *Event) error
	DeleteEvent(id int) error
}
