package bunt

import (
	"testing"
	"time"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/google/go-cmp/cmp"
)

var testDebts = []*model.Debt{
	{0, 271828, 15, "Bernard", "@ck", "@noah", time.Now(), "USD"},
	{1, 271828, 25, "ChiniseCuisine", "@chapp", "@rock", time.Now(), "USD"},
	{2, 271828, 13, "TouchMyBao", "@murphy", "@chapp", time.Now(), "USD"},
	{3, 271828, 13, "TouchMyBao", "@murphy", "@rock", time.Now(), "USD"},
	{4, 271828, 13, "TouchMyBao", "@murphy", "@noah", time.Now(), "USD"},
}

func testDebtRepository(t *testing.T) model.DebtRepository {
	repo, err := NewDebtRepository(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func TestStoreFind(t *testing.T) {
	repo := testDebtRepository(t)
	for _, debt := range testDebts {
		err := repo.Store(debt)
		failOnError(t, err, "failed to store debt")
	}

	expectedDebt := testDebts[len(testDebts)/2]
	gotDebt, err := repo.Find(expectedDebt.GroupID, expectedDebt.ID)
	failOnError(t, err, "failed to find debt")
	if diff := cmp.Diff(expectedDebt, gotDebt); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}

func TestStoreFindAll(t *testing.T) {
	repo := testDebtRepository(t)
	for _, debt := range testDebts {
		repo.Store(debt)
	}
	all, err := repo.FindAll(testDebts[0].GroupID)
	failOnError(t, err, "failed to find all debts")
	if len(all) != len(testDebts) {
		t.Fatalf("expected %d debts but got %d", len(testDebts), len(all))
	}
	for _, storedDebt := range all {
		id := storedDebt.ID
		if diff := cmp.Diff(testDebts[id], storedDebt); diff != "" {
			t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
		}
	}
}

func failOnError(t *testing.T, err error, context string) {
	if err != nil {
		t.Fatalf("%s: %v", context, err)
	}
}
