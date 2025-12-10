package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Vote struct {
	ID          int
	ElectionID  int
	CandidateID int
	UserID      int
	CreatedAt   time.Time
}

type VoteService struct {
	DB *sql.DB
}

func (v *VoteService) AddVote(userId, candidateId, electionId int) (*Vote, error) {
	vote := Vote{
		CandidateID: candidateId,
		ElectionID:  electionId,
		UserID:      userId,
		CreatedAt:   time.Now(),
	}

	row := v.DB.QueryRow(`INSERT INTO votes (election_id, candidate_id, user_id, created_id) VALUES ($1, $2, $3) RETURNING id`,
		vote.ElectionID, vote.CandidateID, vote.UserID, vote.CreatedAt)

	err := row.Scan(&vote.ID)
	if err != nil {
		return nil, fmt.Errorf("voting: %v", err)
	}

	return &vote, nil
}

