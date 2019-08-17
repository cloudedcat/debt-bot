package manager

import (
	"testing"

	"github.com/cloudedcat/debt-bot/bunt"
	"github.com/cloudedcat/debt-bot/model"
	"github.com/cloudedcat/debt-bot/testset"
	"github.com/tidwall/buntdb"
)

func testOpen(t *testing.T) *buntdb.DB {
	db, err := bunt.Open(":memory:")
	testset.FatalOnError(t, err, "failed to open db connection")
	return db
}

func TestRegisterGroup(t *testing.T) {
	db := testOpen(t)
	groups, partics := bunt.NewGroupRepository(db), bunt.NewParticipantRepository(db)
	service := NewService(groups, partics)

	err := service.RegisterGroup(model.Group{ID: testset.GroupID})

	testset.FatalOnError(t, err, "faied to register group")
	_, err = groups.Find(testset.GroupID)
	testset.FatalOnError(t, err, "faied to find registered group")
}

func TestRegisterParticipants(t *testing.T) {
	db := testOpen(t)
	groups, partics := bunt.NewGroupRepository(db), bunt.NewParticipantRepository(db)
	service := NewService(groups, partics)
	err := service.RegisterGroup(model.Group{ID: testset.GroupID})
	testset.FatalOnError(t, err, "faied to register group")

	for _, partic := range testset.Participants {
		err := service.RegisterParticipant(testset.GroupID, *partic)
		testset.FatalOnError(t, err, "faied to register participant")
	}

	expected := testset.Participants[len(testset.Participants)/2]
	_, err = partics.Find(testset.GroupID, expected.ID)
	testset.FatalOnError(t, err, "faied to find participant")
}
