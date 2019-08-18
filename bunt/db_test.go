package bunt

import (
	"os"
	"testing"

	"github.com/cloudedcat/debt-bot/model"
	"github.com/cloudedcat/debt-bot/testset"
	"github.com/google/go-cmp/cmp"
)

func TestIndexRestoring(t *testing.T) {
	persDB := "persistent.db"
	db, err := Open(persDB)
	defer os.Remove(persDB)

	testset.FatalOnError(t, err, "failed to open db")
	testGroupID := model.GroupID(-260219127)
	repo := NewGroupRepository(db)
	repo.Store(&model.Group{ID: testGroupID})
	expectedIndexes, err := db.Indexes()
	testset.FatalOnError(t, err, "failed to get indexes")

	err = db.Close()
	testset.FatalOnError(t, err, "failed to close db")

	db, err = Open(persDB)
	testset.FatalOnError(t, err, "failed to reopen db")

	gotIndexes, err := db.Indexes()
	testset.FatalOnError(t, err, "failed to get indexes")

	if diff := cmp.Diff(expectedIndexes, gotIndexes); diff != "" {
		t.Fatalf("Indexes mismatch (-expected, +got):\n%s", diff)
	}
}
