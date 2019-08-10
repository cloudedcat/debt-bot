package calculator

import (
	"testing"

	"github.com/cloudedcat/finance-bot/bunt"
	"github.com/cloudedcat/finance-bot/model"
	"github.com/cloudedcat/finance-bot/testset"
	"github.com/tidwall/buntdb"
)

func testOpen(t *testing.T) *buntdb.DB {
	db, err := bunt.Open(":memory:")
	testset.FatalOnError(t, err, "failed to open db connection")
	return db
}

func addGroup(t *testing.T, db *buntdb.DB) {
	groups := bunt.NewGroupRepository(db)
	newGroup := model.BuildGroup(testset.GroupID)
	err := groups.Store(newGroup)
	testset.FatalOnError(t, err, "failed to store group")
}

func addParticipants(t *testing.T, db *buntdb.DB) {
	partics := bunt.NewParticipantRepository(db)
	for _, p := range testset.Participants {
		err := partics.Store(testset.GroupID, p)
		testset.FatalOnError(t, err, "failed to store participant")
	}
}

func addDebtsViaService(t *testing.T, service Service) {
	for _, d := range testset.Debts {
		err := service.AddDebt(testset.GroupID, *d)
		testset.FatalOnError(t, err, "failed to add debt via service")
	}
}

func TestAddDebt(t *testing.T) {
	db := testOpen(t)
	addGroup(t, db)
	addParticipants(t, db)
	debts, partics := bunt.NewDebtRepository(db), bunt.NewParticipantRepository(db)
	service := NewService(debts, partics)
	addDebtsViaService(t, service)

	got, err := debts.FindAll(testset.GroupID)
	testset.FatalOnError(t, err, "failed to find all debts")
	if len(got) != len(testset.Debts) {
		t.Fatalf("expected %d debts but got %d", len(testset.Debts), len(got))
	}
}

func TestCalculateDebts(t *testing.T) {
	db := testOpen(t)
	addGroup(t, db)
	addParticipants(t, db)
	debts, partics := bunt.NewDebtRepository(db), bunt.NewParticipantRepository(db)
	service := NewService(debts, partics)
	addDebtsViaService(t, service)

	fDebts, err := service.CalculateDebts(testset.GroupID)
	testset.FatalOnError(t, err, "failed to calculate debts")
	for _, fDebt := range fDebts {
		if fDebt.Borrower.Alias == "" || fDebt.Lender.Alias == "" || fDebt.Debt.Amount == 0 {
			t.Fatal("Final debt has an empty Alias or zero Amount")
		}
	}
}
