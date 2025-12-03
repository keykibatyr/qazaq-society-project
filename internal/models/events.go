package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Event struct {
	ID          int
	Title       string
	Description string
	ImageURL    string

	StartDate time.Time

	Published bool

	Day   string
	Month string
	Time  string
}

type EventService struct {
	DB *sql.DB
}

func (es *EventService) CreateEvent(title, description, image string, start time.Time, published bool) (*Event, error) {
	newEvent := Event{
		Title:       title,
		Description: description,
		ImageURL:    image,
		StartDate:   start,
		Published:   published,
	}

	newEvent.Day = newEvent.StartDate.Format("02")
	newEvent.Month = newEvent.StartDate.Month().String()
	newEvent.Time = newEvent.StartDate.Format("3:04 PM")



	row := es.DB.QueryRow(`INSERT INTO events (title, description, image_URL,
	 start_date, published) VALUES ($1, $2, $3, $4, $5) RETURNING id`, newEvent.Title, newEvent.Description,
		newEvent.ImageURL, newEvent.StartDate, newEvent.Published)

	fmt.Println("bug1")

	err := row.Scan(&newEvent.ID)
	if err != nil {
		return nil, fmt.Errorf("could insert: %v", err)
	}

	fmt.Println("bug2")
	return &newEvent, nil
}

func (es *EventService) GetAll() ([]Event, error) {
	rows, err := es.DB.Query(`SELECT id, title, description, image_URL,
	 start_date, published FROM events ORDER BY start_date ASC`)
	if err != nil {
		return nil, fmt.Errorf("select all: %v", err)
	}

	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description,
		&event.ImageURL, &event.StartDate, &event.Published)
		if err != nil {
			return nil, fmt.Errorf("scanning a row: %v", err)
		}

		event.Day = event.StartDate.Format("02")
		event.Month = event.StartDate.Month().String()
		event.Time = event.StartDate.Format("3:04 PM")

		events = append(events, event)
	}

	return events, nil
}

func (es *EventService) GetLatest() (*Event, error) {
	var event Event 
	row := es.DB.QueryRow(`SELECT id, title, description, image_URL,
	 start_date, published FROM events WHERE start_date > NOW() 
	 ORDER BY start_date ASC LIMIT 1;`)
	err := row.Scan(&event.ID, &event.Title, &event.Description,
		&event.ImageURL, &event.StartDate, &event.Published)
	if err != nil {
		return nil, fmt.Errorf("scanning a row: %v", err)
	}

	event.Day = event.StartDate.Format("02")
	event.Month = event.StartDate.Month().String()
	event.Time = event.StartDate.Format("3:04 PM")

	return &event, nil
}

func (es *EventService) Publish(id int) error {
	_, err := es.DB.Exec(`UPDATE events SET published = TRUE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publishing: %v", err)
	}

	return nil
}


func (es *EventService) UnPublish(id int) error {
	_, err := es.DB.Exec(`UPDATE events SET published = FALSE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publishing: %v", err)
	}

	return nil
}

func (es *EventService) Delete(id int) error {
	_, err := es.DB.Exec(`DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publishing: %v", err)
	}

	return nil
}
