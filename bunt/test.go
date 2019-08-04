package bunt

import (
	"testing"

	"github.com/tidwall/buntdb"
)

func failOnError(t *testing.T, err error, context string) {
	if err != nil {
		t.Fatalf("%s: %v", context, err)
	}
}

func testOpen(t *testing.T) *buntdb.DB {
	db, err := Open(":memory:")
	failOnError(t, err, "failed to open db connection")
	return db
}
