package bunt

import (
	"fmt"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/tidwall/buntdb"
)

const indexGroup = "groups"
const prefixGroup = "debt::"

type groupRepository struct {
	db *buntdb.DB
}

// NewGroupRepository returns new instance of a BuntDB group repository
func NewGroupRepository(dbName string) (model.GroupRepository, error) {
	db, err := buntdb.Open(dbName)
	if err != nil {
		return nil, err
	}

	return &groupRepository{db: db}, nil
}

func (r *groupRepository) Find(id model.GroupID) (*model.Group, error) {
	var val string
	err := r.db.View(func(tx *buntdb.Tx) error {
		var err error
		val, err = tx.Get(fmt.Sprintf("%s%d", prefixGroup, id))
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

func parseGroup(raw string) (*model.Group, error) {
	return nil, nil
}
