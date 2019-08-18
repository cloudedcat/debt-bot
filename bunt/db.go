package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/debt-bot/model"
	"github.com/tidwall/buntdb"
)

const sep = ":"

// groupSet contains all groups
var groupSet []model.GroupID

// Open opens a database at the provided path.
// If the database does not exist then it will be created automatically.
func Open(path string, shouldRestoreIndex bool) (*buntdb.DB, error) {
	// dbExists := doesDBExist(path)

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	// restore all indexes because Bunt doesn't save them
	err = db.Update(func(tx *buntdb.Tx) error {
		tx.CreateIndex(indexGroup, patternByIndex(indexGroup), buntdb.IndexJSON("ID"))
		if err := restoreGroupSet(tx); err != nil {
			return err
		}
		if shouldRestoreIndex {
			return restoreIndexes(tx)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}

func restoreGroupSet(tx *buntdb.Tx) error {
	val, err := tx.Get("group_set")
	if err == buntdb.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(val), &groupSet); err != nil {
		return err
	}
	return nil
}

func restoreIndexes(tx *buntdb.Tx) error {
	for _, gID := range groupSet {
		if err := createDebtParticipantIndexes(tx, gID); err != nil {
			return err
		}
	}
	return nil
}

func createDebtParticipantIndexes(tx *buntdb.Tx, id model.GroupID) error {
	indexes := []string{indexDebt(id), indexParticipant(id)}
	for _, index := range indexes {
		err := tx.CreateIndex(index, patternByIndex(index), buntdb.IndexJSON("ID"))
		if err != nil {
			return err
		}
	}
	return nil
}

func patternByIndex(index string) string {
	return fmt.Sprintf("%s%s*", index, sep)
}
