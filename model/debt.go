package model

import "time"

// Currency is unused for now type
type Currency string

// Debt represents a debt. The main entity in the domain.
type Debt struct {
	ID       int
	GroupID  GroupID
	Amount   float64
	Tag      string
	Borrower ParticipantID
	Lender   ParticipantID
	Date     time.Time
	// Reserved fields for future purposes
	Currency Currency
}

// func New() {

// }

// DebtRepository provides access to a debt store.
type DebtRepository interface {
	FindAll(groupID GroupID) ([]*Debt, error)
	Find(groupID GroupID, id int) (*Debt, error)
	Store(debt *Debt) error
	// NextID(chatID ChatID) (int, error)
}
