package bunt

import (
	"fmt"
	"testing"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/buntdb"
)

func testParticipantRepository(t *testing.T, conn *buntdb.DB) model.ParticipantRepository {
	return NewParticipantRepository(conn)
}

var testParticipants = []*model.Participant{
	{12, "@ck", "Louis", "C K"},
	{13, "@noah", "Trevor", "Noah"},
	{14, "@chapp", "Dave", "Chappelle"},
	{15, "@murphy", "Eddie", "Murphy"},
	{16, "@rock", "Chris", "Rock"},
}

func uploadTestParticipants(t *testing.T, groupID model.GroupID, repo model.ParticipantRepository) {
	for _, partic := range testParticipants {
		err := repo.Store(groupID, partic)
		failOnError(t, err, fmt.Sprintf("failed to load participant: %v", partic.ID))
	}
}

func TestParticipantStoreFind(t *testing.T) {
	db := testOpen(t)
	gRepo, repo := NewGroupRepository(db), NewParticipantRepository(db)
	uploadTestGroup(t, model.BuildGroup(testGroupID), gRepo)
	uploadTestParticipants(t, testGroupID, repo)
	expected := testParticipants[len(testParticipants)/2]

	got, err := repo.Find(testGroupID, expected.ID)

	failOnError(t, err, "repository error")
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}

func TestParticipantFindAll(t *testing.T) {
	db := testOpen(t)
	gRepo, repo := NewGroupRepository(db), NewParticipantRepository(db)
	uploadTestGroup(t, model.BuildGroup(testGroupID), gRepo)
	uploadTestParticipants(t, testGroupID, repo)
	expected := testParticipants

	got, err := repo.FindAll(testGroupID)
	failOnError(t, err, "repository error")
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}
