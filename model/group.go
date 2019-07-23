package model

type GroupID int

type Group struct {
	ID           GroupID
	// Participants []Participant
	// Debts        []Debt
}

// GroupRepository provides access to a group store.
type GroupRepository interface {
	// FindAll(groupID GroupID) ([]*Participant, error)
	Find(groupID GroupID) (*Group, error)
	Store(group *Group) error
	// NextID(chatID ChatID) (int, error)
}
