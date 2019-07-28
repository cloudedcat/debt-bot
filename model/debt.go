package model

import "time"

// Currency - NYI
type Currency string

// DebtID is a custom type for identifying DebtID entity.
// ID generating should be provided debt repository.
type DebtID int

// Debt represents a debt. The main entity in the domain.
type Debt struct {
	ID       int
	GroupID  GroupID
	Amount   float64
	Tag      string
	Borrower ParticipantID
	Lender   ParticipantID
	Date     time.Time
	// Fields reserved for future purposes
	Currency Currency
}

// DebtRepository provides access to a debt store.
type DebtRepository interface {
	FindAll(groupID GroupID) ([]*Debt, error)
	Find(groupID GroupID, id int) (*Debt, error)
	Store(debt *Debt) error
	NextID(groupID GroupID) (DebtID, error)
}
