package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/tidwall/buntdb"
)

const prefixParticipant = "participant"

func indexParticipant(groupID model.GroupID) string {
	return fmt.Sprintf("%s%s%d", prefixParticipant, sep, int(groupID))
}

func patternParticipant(groupID model.GroupID) string {
	return indexParticipant(groupID)
}

// NewParticipantRepository creates new instanse of ParticipantRepository
func NewParticipantRepository(db *buntdb.DB) model.ParticipantRepository {
	return &participantRepository{db: db}
}

type participantRepository struct {
	db *buntdb.DB
}

func (p *participantRepository) Find(
	groupID model.GroupID, id model.ParticipantID) (*model.Participant, error) {

	var raw string
	var err error
	err = p.db.View(func(tx *buntdb.Tx) error {
		raw, err = tx.Get(p.key(groupID, id))
		return err
	})
	if err != nil {
		return nil, err
	}

	return parseParticipant(raw)
}

func (p *participantRepository) Store(groupID model.GroupID, partic *model.Participant) error {
	return p.db.Update(func(tx *buntdb.Tx) error {
		composedPartic, err := composeParticipant(partic)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(p.key(groupID, partic.ID), composedPartic, nil)
		return err
	})
}

func parseParticipant(raw string) (*model.Participant, error) {
	p := &model.Participant{}
	if err := json.Unmarshal([]byte(raw), p); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *participantRepository) key(groupID model.GroupID, id model.ParticipantID) string {
	return fmt.Sprintf("%s%s%d", patternParticipant(groupID), sep, int(id))
}

func composeParticipant(p *model.Participant) (string, error) {
	b, err := json.Marshal(p)
	return string(b), err
}
