package model

import (
	"errors"
	"time"
)

// Currency - NYI
type Currency string

// DebtID is a custom type for identifying DebtID entity.
// ID generating should be provided debt repository.
type DebtID int

// Debt represents a debt. The main entity in the domain.
type Debt struct {
	ID         DebtID
	Amount     float64
	Tag        string
	BorrowerID ParticipantID
	LenderID   ParticipantID
	Date       time.Time
	// Fields reserved for future purposes
	Currency Currency
}

var ErrBlankField = errors.New("blank field")
var ErrParticipantCollision = errors.New("borrower cant be lender")

func (d *Debt) Validate() error {

	if d.Amount == 0 {
		return ErrBlankField
	}
	if d.BorrowerID == d.LenderID {
		return ErrParticipantCollision
	}
	return nil
}

// DebtRepository provides access to a debt store.
type DebtRepository interface {
	FindAll(groupID GroupID) ([]*Debt, error)
	Find(groupID GroupID, id DebtID) (*Debt, error)
	Store(groupID GroupID, debts ...*Debt) error
	NextID(groupID GroupID) (DebtID, error)
	Clear(groupID GroupID) error
}
