package model

// GroupID is a custom type for identifying Group entity.
// It assigns Telegram ChatID to GroupID
type GroupID int64

// Group represents a Telegram Chat.
// Participant and Debt models link to Group and have no meaning without it
type Group struct {
	ID GroupID
}

// BuildGroup is used to create a Group instance
func BuildGroup(id GroupID) *Group {
	return &Group{ID: id}
}

func (g *Group) Validate() error {
	return nil
}

// GroupRepository provides access to a group store.
type GroupRepository interface {
	Find(groupID GroupID) (*Group, error)
	Store(group *Group) error
	// NextID(chatID ChatID) (int, error)
}
