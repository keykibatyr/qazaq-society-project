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

type CandidateVoteCount struct {
	CandidateID int
	CandidateName string
	VoteCount int
}

func(v *VoteService) VotePerCandidate(electionID int) ([]CandidateVoteCount, error) {
	var candidateVoteCount []CandidateVoteCount

	rows, err := v.DB.Query(`SELECT COUNT(v.candidate_id), c.id
	 FROM candidates c
	 LEFT JOIN votes v
	 ON v.election_id = c.election_id AND
	 v.candidate_id = c.id WHERE c.election_id = $1
	 GROUP BY c.id`, electionID)
	if err != nil {
		return nil, fmt.Errorf("extracting coandidate-votes: %v", err)
	}

	for rows.Next() {
		var candidateVC CandidateVoteCount

		err := rows.Scan(&candidateVC.VoteCount, &candidateVC.CandidateID)

		if err != nil {
			return nil, fmt.Errorf("extracting candidate-votes row: %v", err)
		}

		candidateVoteCount = append(candidateVoteCount, candidateVC)
	
	}

	fmt.Println(candidateVoteCount)

	return candidateVoteCount, nil 

	//REFACTOOOOR ADD ROWS WITH CANDIDATES WITH 0 VOTES TOO 
}
