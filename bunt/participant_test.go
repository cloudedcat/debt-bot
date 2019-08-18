package bunt

import (
	"fmt"
	"testing"

	"github.com/cloudedcat/debt-bot/model"
	"github.com/cloudedcat/debt-bot/testset"
	"github.com/google/go-cmp/cmp"
)

func uploadTestParticipants(t *testing.T, groupID model.GroupID, repo model.ParticipantRepository) {
	for _, partic := range testset.Participants {
		err := repo.Store(groupID, partic)
		testset.FatalOnError(t, err, fmt.Sprintf("failed to load participant: %v", partic.ID))
	}
}

func TestParticipantStoreFind(t *testing.T) {
	db := testOpen(t)
	defer db.Close()

	gRepo, repo := NewGroupRepository(db), NewParticipantRepository(db)
	uploadTestGroup(t, model.BuildGroup(testset.GroupID), gRepo)
	uploadTestParticipants(t, testset.GroupID, repo)
	expected := testset.Participants[len(testset.Participants)/2]

	got, err := repo.Find(testset.GroupID, expected.ID)

	testset.FatalOnError(t, err, "repository error")
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}

func TestParticipantFindAll(t *testing.T) {
	db := testOpen(t)
	defer db.Close()

	gRepo, repo := NewGroupRepository(db), NewParticipantRepository(db)
	uploadTestGroup(t, model.BuildGroup(testset.GroupID), gRepo)
	uploadTestParticipants(t, testset.GroupID, repo)
	expected := model.Participants(testset.Participants)

	got, err := repo.FindAll(testset.GroupID)
	testset.FatalOnError(t, err, "repository error")
	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Debt mismatch (-expected, +got):\n%s", diff)
	}
}
