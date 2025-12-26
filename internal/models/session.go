package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/keykibatyr/qazaq-society-project/rand"
)

const (
	MinbytesperToken = 32
	TimeSession = 24 * time.Hour
)

type Session struct {
	ID        int
	UserID    int
	Token string
	TokenHash string
	ExpiresAt time.Time
}

type SessionService struct {
	DB *sql.DB

	BytesPerToken int
}

func (s *SessionService) Create(userID int) (*Session, error) {
	bytesPerToken := s.BytesPerToken
	if bytesPerToken < MinbytesperToken {
		bytesPerToken = MinbytesperToken
	}
	
	token, err := rand.ToString(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("token: %v", err)
	}


	tokenHash := s.hash(token)

	session := Session{
		UserID: userID,
		Token: token,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(TimeSession),
	}

	row := s.DB.QueryRow(`INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)
	 ON CONFLICT (user_id) DO UPDATE SET token_hash = $2, expires_at = $3 RETURNING id`, session.UserID, session.TokenHash, session.ExpiresAt)


	err = row.Scan(&session.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting into table: %v", err)
	}

	return &session, nil
}

func (s *SessionService) User(token string) (*User, error) {
	tokenHash := s.hash(token)

	var User User
	row := s.DB.QueryRow(`SELECT users.id, users.email, users.password_hash, users.first_name,
	 users.last_name, users.can_vote, users.created_at, users.role FROM sessions INNER JOIN users ON users.id = sessions.user_id WHERE sessions.token_hash = $1`, tokenHash)
	
	err := row.Scan(&User.ID, &User.Email, &User.PasswordHash,
		 &User.FirstName, &User.LastName, 
		 &User.CanVote, &User.CreatedAt, &User.Role)

	fmt.Printf("this is function 2: %+v\n", User.ID)
	
	if err != nil {
		return nil, fmt.Errorf("getting User: %v", err)
	}

	return &User, nil
		
}

func (s *SessionService) Delete(token string) error {
	tokenHash := s.hash(token)

	_, err := s.DB.Exec(`DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("deleting: %v", err)
	}

	return nil
}

func (s *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (s *SessionService) IsExpired(userID int) (bool, error) { //returns true if expired 
	var expireDate time.Time

	row := s.DB.QueryRow(`SELECT expires_at, token_hash FROM sessions WHERE user_id = $1`, userID)
	err := row.Scan(&expireDate)
	if err != nil {
		return false, fmt.Errorf("could not scan for the expiration date")
	}

	return expireDate.Before(time.Now()), nil
}