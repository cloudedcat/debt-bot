package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/tidwall/buntdb"
)

const indexGroup = "group"
const prefixGroup = indexGroup + sep

type groupRepository struct {
	db *buntdb.DB
}

// NewGroupRepository returns new instance of a BuntDB group repository
func NewGroupRepository(db *buntdb.DB) model.GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) key(id model.GroupID) string {
	return fmt.Sprintf("%s%d", prefixGroup, id)
}

func (r *groupRepository) Find(id model.GroupID) (*model.Group, error) {
	var val string
	err := r.db.View(func(tx *buntdb.Tx) error {
		var err error
		val, err = tx.Get(r.key(id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return parseGroup(val)
}

func (r *groupRepository) Store(group *model.Group) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		indexes := []string{indexDebt(group.ID), indexParticipant(group.ID)}
		for _, index := range indexes {
			err := tx.CreateIndex(index, patternByIndex(index), buntdb.IndexJSON("ID"))
			if err != nil {
				return err
			}
		}
		composedGroup, err := composeGroup(group)
		if err != nil {
			return err
		}
		tx.Set(r.key(group.ID), composedGroup, nil)
		return nil
	})
}

func parseGroup(raw string) (*model.Group, error) {
	group := &model.Group{}
	if err := json.Unmarshal([]byte(raw), group); err != nil {
		return nil, err
	}
	return group, nil
}

func composeGroup(group *model.Group) (string, error) {
	bGroup, err := json.Marshal(group)
	return string(bGroup), err
}
