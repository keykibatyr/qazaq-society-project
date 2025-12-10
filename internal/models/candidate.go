package models

import (
	"database/sql"
	"fmt"
)

type Candidate struct {
	ID         int
	ElectionID int
	Name       string
	ImageURL   string
}

type CandidateService struct {
	DB *sql.DB
}

func (c *CandidateService) CreateCandidate(name, image string, electionId int) (*Candidate, error) {
	candidate := Candidate{
		ElectionID: electionId,
		Name:       name,
		ImageURL:   image,
	}

	row := c.DB.QueryRow(`INSERT INTO candidates (name, image_url, election_id) VALUES ($1, $2, $3) RETURNING id`, candidate.Name, candidate.ImageURL, candidate.ElectionID)
	err := row.Scan(&candidate.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting candidate: %v", err)
	}

	return &candidate, nil

}

func (c CandidateService) GetAllCandidates(electionId int) ([]Candidate, error) {
	var candidates []Candidate

	row, err := c.DB.Query(`SELECT id, name, election_id, image_url FROM candidates WHERE election_id = $1`, electionId)
	if err != nil {
		return nil, fmt.Errorf("getting all candidates: %v", err)
	}

	for row.Next(){
		var candidate Candidate
		err = row.Scan(&candidate.ID, &candidate.Name, &candidate.ElectionID, &candidate.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("getting a candidate: %v", err)
		}

		candidates = append(candidates, candidate)
	}

	return candidates, nil
}