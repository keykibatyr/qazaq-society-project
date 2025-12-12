package models

import (
	"database/sql"
	"errors"
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

	row := v.DB.QueryRow(`INSERT INTO votes (election_id, candidate_id, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id`,
		vote.ElectionID, vote.CandidateID, vote.UserID, vote.CreatedAt)

	err := row.Scan(&vote.ID)
	if err != nil {
		return nil, fmt.Errorf("voting: %v", err)
	}

	return &vote, nil
}

func (v *VoteService) VoteCheck(userId, electionId int) bool {
	row := v.DB.QueryRow(`SELECT id FROM votes WHERE user_id = $1 AND election_id = $2`, userId, electionId)

	var voteExists int

	err := row.Scan(&voteExists)
	if errors.Is(err, sql.ErrNoRows) {
		return true
	}

	if err != nil {
		return false
	}

	return false
}
