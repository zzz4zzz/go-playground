package models

import (
	"time"
)

type ModelTimer interface {
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m BaseModel) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m BaseModel) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}

type User struct {
	BaseModel
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	BankAccounts []BankAccount `gorm:"many2many:user_bank_accounts" json:"bank_accounts"`
}

type BankAccount struct {
	BaseModel
	Name   string `json:"name"`
	Number string `json:"number"`
}
