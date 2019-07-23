package model

type ParticipantID int
type Alias string

// Participant represents person who take part in sharing debts
type Participant struct {
	ID ParticipantID
	Alias      Alias
	FirstName  string
	LastName   string
}


// ParticipantRepository provides access to a participant store.
type ParticipantRepository interface {
	// FindAll(groupID GroupID) ([]*Participant, error)
	Find(groupID GroupID, id ParticipantID) (*Participant, error)
	FindByAlias(groupID GroupID, alias Alias) (*Participant, error)
	Store(participant *Participant) error
	// NextID(chatID ChatID) (int, error)
}
