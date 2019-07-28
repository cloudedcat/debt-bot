package bunt

import (
	"os"

	"github.com/tidwall/buntdb"
)

const sep = ":"

// Open opens a database at the provided path.
// If the database does not exist then it will be created automatically.
func Open(path string) (*buntdb.DB, error) {
	dbExists := doesDBExist(path)

	db, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	if dbExists {
		return db, nil
	}
	// If the database is just created then add necessary indexes
	err = db.Update(func(tx *buntdb.Tx) error {
		return tx.CreateIndex(indexGroup, prefixGroup)
	})

	if err != nil {
		return nil, err
	}
	return db, nil
}

func doesDBExist(path string) bool {
	inMemory := path == ":memory:"
	if inMemory {
		return false
	}

	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
