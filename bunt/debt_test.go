package bunt

import (
	"fmt"
	"testing"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/cloudedcat/finance-bot/testset"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/buntdb"
)

func testOpen(t *testing.T) *buntdb.DB {
	db, err := Open(":memory:")
	testset.FatalOnError(t, err, "failed to open db connection")
	return db
}

func uploadTestDebts(t *testing.T, groupID model.GroupID, repo model.DebtRepository) {
	var err error
	for _, debt := range testset.Debts {
		debt.ID, err = repo.NextID(groupID)
		testset.FatalOnError(t, err, fmt.Sprintf("failed to get nex id"))
	}

	err = repo.Store(groupID, testset.Debts...)
	testset.FatalOnError(t, err, fmt.Sprintf("failed to upload debt list '%v'", testset.Debts))
}

func testUploadAll(t *testing.T, db *buntdb.DB) {
	gRepo := NewGroupRepository(db)
	pRepo := NewParticipantRepository(db)
	dRepo := NewDebtRepository(db)

	uploadTestGroup(t, model.BuildGroup(testset.GroupID), gRepo)
	uploadTestParticipants(t, testset.GroupID, pRepo)
	uploadTestDebts(t, testset.GroupID, dRepo)
}

func TestDebtStoreFind(t *testing.T) {
	db := testOpen(t)
	testUploadAll(t, db)
	repo := NewDebtRepository(db)

	expectedDebt := testset.Debts[len(testset.Debts)/2]
	gotDebt, err := repo.Find(testset.GroupID, expectedDebt.ID)
	testset.FatalOnError(t, err, "failed to find debt")
	if diff := cmp.Diff(expectedDebt, gotDebt); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}

func TestDebtFindAll(t *testing.T) {
	db := testOpen(t)
	testUploadAll(t, db)
	repo := NewDebtRepository(db)

	got, err := repo.FindAll(testset.GroupID)
	testset.FatalOnError(t, err, "failed to find all debts")
	if diff := cmp.Diff(got, testset.Debts); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}

func TestDebtNextID(t *testing.T) {
	db := testOpen(t)
	testUploadAll(t, db)
	repo := NewDebtRepository(db)
	// counter starts with 0, so next ID equals number of debts
	expected := model.DebtID(len(testset.Debts))

	got, err := repo.NextID(testset.GroupID)
	testset.FatalOnError(t, err, "failed to get next debt ID")
	if got != expected {
		t.Fatalf("expected next id is %v but got %v", expected, got)
	}
}
