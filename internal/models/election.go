package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Election struct {
	ID          int
	Title       string
	Description string
	StartDate   time.Time
	EndDate     time.Time
	Published   bool
	CreatedAt   time.Time

	StartDay   string
	StartMonth string
	StartTime  string

	EndDay   string
	EndMonth string
	EndTime  string
}

type ElectionService struct {
	DB *sql.DB
}

func (e *ElectionService) Create(title, description string, startDate, endDate time.Time, published bool) (*Election, error) {
	election := Election{
		Title:       title,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Published:   published,
		CreatedAt:   time.Now(),
	}

	election.StartDay = election.StartDate.Format("02")
	election.StartMonth = election.StartDate.Month().String()
	election.StartTime = election.StartDate.Format("15:04 PM")

	election.EndDay = election.EndDate.Format("02")
	election.EndMonth = election.EndDate.Month().String()
	election.EndTime = election.EndDate.Format("15:04 PM")

	row := e.DB.QueryRow(`INSERT INTO elections (title, description, start_at, end_at, published, created_at)
	 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, election.Title, election.Description, election.StartDate, election.EndDate, election.Published, election.CreatedAt)

	err := row.Scan(&election.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting into db: %v", err)
	}

	return &election, nil

}

func (e *ElectionService) GetAllElections() ([]Election, error) {
	var elections []Election

	rows, err := e.DB.Query(`SELECT id, title, description, start_at, end_at, published, created_at FROM elections`)
	if err != nil {
		return nil, fmt.Errorf("extracting all elections: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var election Election

		err = rows.Scan(&election.ID, &election.Title, &election.Description,
			 &election.StartDate, &election.EndDate,
			 &election.Published, &election.CreatedAt)

		

		if err != nil {
			return nil, fmt.Errorf("extracting election: %v", err)	
		}

		election.StartDay = election.StartDate.Format("02")
		election.StartMonth = election.StartDate.Month().String()
		election.StartTime = election.StartDate.Format("15:04 PM")

		election.EndDay = election.EndDate.Format("02")
		election.EndMonth = election.EndDate.Month().String()
		election.EndTime = election.EndDate.Format("15:04 PM")


		elections = append(elections, election)
	}

	return elections, nil
}

func (e *ElectionService) UnPublish(id int) error {
	_, err := e.DB.Exec(`UPDATE elections SET published = FALSE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publishing: %v", err)
	}

	return nil
}

func (e *ElectionService) Publish(id int) error {
	_, err := e.DB.Exec(`UPDATE elections SET published = TRUE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("publishing: %v", err)
	}

	return nil
}