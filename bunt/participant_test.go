package bunt

import (
	"testing"

	"github.com/cloudedcat/finance-bot/model"
)

func testParticipantRepository(t *testing.T) model.ParticipantRepository {
	db, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return NewParticipantRepository(db)
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
		if err := repo.Store(groupID, partic); err != nil {
			t.Fatal(err)
		}
	}
}

func TestParticipantStoreFind(t *testing.T) {
	repo := testParticipantRepository(t)
	uploadTestParticipants(t, testGroupID, repo)
}
