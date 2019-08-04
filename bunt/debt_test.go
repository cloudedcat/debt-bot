package bunt

import (
	"fmt"
	"testing"
	"time"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/buntdb"
)

const testCurrency = "USD"

var testDebts = []*model.Debt{
	{0, 15, "Bernard", testParticipants[0].ID, testParticipants[1].ID, time.Now(), testCurrency},
	{1, 25, "ChiniseCuisine", testParticipants[1].ID, testParticipants[2].ID, time.Now(), testCurrency},
	{2, 13, "BaoStore", testParticipants[3].ID, testParticipants[2].ID, time.Now(), testCurrency},
	{3, 13, "BaoStore", testParticipants[3].ID, testParticipants[1].ID, time.Now(), testCurrency},
	{4, 13, "BaoStore", testParticipants[2].ID, testParticipants[0].ID, time.Now(), testCurrency},
}

func uploadTestDebts(t *testing.T, groupID model.GroupID, repo model.DebtRepository) {
	for _, debt := range testDebts {
		err := repo.Store(groupID, debt)
		failOnError(t, err, fmt.Sprintf("failed to load debt: %v", debt.ID))
	}
}

func testUploadAll(t *testing.T, db *buntdb.DB) {
	gRepo := NewGroupRepository(db)
	pRepo := NewParticipantRepository(db)
	dRepo := NewDebtRepository(db)

	uploadTestGroup(t, model.BuildGroup(testGroupID), gRepo)
	uploadTestParticipants(t, testGroupID, pRepo)
	uploadTestDebts(t, testGroupID, dRepo)
}

func TestStoreFind(t *testing.T) {
	db := testOpen(t)
	testUploadAll(t, db)
	repo := NewDebtRepository(db)

	expectedDebt := testDebts[len(testDebts)/2]
	gotDebt, err := repo.Find(testGroupID, expectedDebt.ID)
	failOnError(t, err, "failed to find debt")
	if diff := cmp.Diff(expectedDebt, gotDebt); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}

func TestStoreFindAll(t *testing.T) {
	db := testOpen(t)
	testUploadAll(t, db)
	repo := NewDebtRepository(db)

	got, err := repo.FindAll(testGroupID)
	failOnError(t, err, "failed to find all debts")
	if diff := cmp.Diff(got, testDebts); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}
