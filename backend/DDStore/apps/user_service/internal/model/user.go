package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"create_at"`
	UpdatedAt      time.Time `json:"update_at"`
	DeletedAt      time.Time `json:"delete_at"`
}

type DeleteUser struct {
	Id        int       `json:"id"`
	Delete_at time.Time `json:"delete_at"`
}

func (User) GetTableName() string { return "users" }

func (user User) VerifyPassword(pw []byte) error {
	err := bcrypt.CompareHashAndPassword(user.HashedPassword, pw)
	return err
}

func (user *User) SetHashPassword(pw []byte) {
	hashedPassword := hashPassword(pw)
	user.HashedPassword = hashedPassword
}

func (user User) Validate() error {
	return nil
}

func hashPassword(pw []byte) []byte {
	hashedPassword, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hashPassword error: %v", err)
	}
	return hashedPassword
}

type Paging struct {
	Limit  int
	Offset int
}

func (p *Paging) Process() {
	if p.Limit > MaxLimit {
		p.Limit = MaxLimit
	}
	if p.Limit < 0 {
		p.Limit = 0
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}
