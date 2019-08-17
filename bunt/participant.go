package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/debt-bot/model"
	"github.com/tidwall/buntdb"
)

func prefixParticipant(groupID model.GroupID) string {
	return fmt.Sprintf("participant%s%d", sep, int(groupID))
}

func indexParticipant(groupID model.GroupID) string {
	return prefixParticipant(groupID)
}

// NewParticipantRepository creates new instance of ParticipantRepository
func NewParticipantRepository(db *buntdb.DB) model.ParticipantRepository {
	return &participantRepository{db: db}
}

type participantRepository struct {
	db *buntdb.DB
}

func (r *participantRepository) Find(
	groupID model.GroupID, id model.ParticipantID) (*model.Participant, error) {

	var raw string
	var err error
	err = r.db.View(func(tx *buntdb.Tx) error {
		raw, err = tx.Get(r.key(groupID, id))
		return err
	})
	if err != nil {
		return nil, err
	}

	return r.parse(raw)
}

func (r *participantRepository) FindAll(groupID model.GroupID) (model.Participants, error) {
	var participants []*model.Participant
	var err = r.db.View(func(tx *buntdb.Tx) error {
		var pErr error
		txErr := tx.Ascend(indexParticipant(groupID), func(_, raw string) bool {
			var parsed *model.Participant
			if parsed, pErr = r.parse(raw); pErr != nil {
				return false
			}
			participants = append(participants, parsed)
			return true
		})
		if pErr != nil {
			return pErr
		}
		return txErr
	})
	if err != nil {
		return nil, err
	}

	return model.Participants(participants), nil
}

func (r *participantRepository) Store(groupID model.GroupID, partic *model.Participant) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		composedPartic, err := r.compose(partic)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(r.key(groupID, partic.ID), composedPartic, nil)
		return err
	})
}

func (r *participantRepository) key(groupID model.GroupID, id model.ParticipantID) string {
	return fmt.Sprintf("%s%s%d", prefixParticipant(groupID), sep, int(id))
}

func (r *participantRepository) parse(raw string) (*model.Participant, error) {
	partic := &model.Participant{}
	if err := json.Unmarshal([]byte(raw), partic); err != nil {
		return nil, err
	}
	return partic, nil
}
func (r *participantRepository) compose(p *model.Participant) (string, error) {
	bRaw, err := json.Marshal(p)
	return string(bRaw), err
}
