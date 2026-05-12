package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Id             int            `json:"id" gorm:"column:id"`
	Name           string         `json:"name" gorm:"column:name"`
	Email          string         `json:"email" gorm:"column:email"`
	HashedPassword []byte         `gorm:"column:password"`
	Role           string         `json:"role" gorm:"column:role"`
	CreatedAt      time.Time      `json:"create_at"`
	UpdatedAt      time.Time      `json:"update_at"`
	DeletedAt      gorm.DeletedAt `json:"delete_at" gorm:"index"`
}

type DeleteUser struct {
	Id        int       `json:"id" gorm:"column:id"`
	Delete_at time.Time `json:"delete_at" gorm:"column:delete_at"`
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
