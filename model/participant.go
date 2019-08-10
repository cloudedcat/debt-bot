package model

// ParticipantID is Telegram ID
type ParticipantID int

// Alias is Telegram nickname started with @.
type Alias string

// Participant represents person who take part in sharing debts.
type Participant struct {
	ID        ParticipantID
	Alias     Alias
	FirstName string
	LastName  string
}

// Participants is an array of *Participant.
type Participants []*Participant

// AsMap converts []*Participant to map with ParticipantID as key
func (ps Participants) AsMap() map[ParticipantID]*Participant {
	mPartics := make(map[ParticipantID]*Participant)
	for _, partic := range ps {
		mPartics[partic.ID] = partic
	}
	return mPartics
}

// ParticipantRepository provides access to a participant store.
type ParticipantRepository interface {
	FindAll(groupID GroupID) (Participants, error)
	Find(groupID GroupID, id ParticipantID) (*Participant, error)
	// FindByAlias(groupID GroupID, alias Alias) (*Participant, error)
	Store(groupID GroupID, participant *Participant) error
}
