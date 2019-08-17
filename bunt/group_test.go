package bunt

import (
	"fmt"
	"testing"

	"github.com/cloudedcat/debt-bot/model"
	"github.com/cloudedcat/debt-bot/testset"
)

func testGroupRepository(t *testing.T) model.GroupRepository {
	repo, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return &groupRepository{db: repo}
}

func uploadTestGroup(t *testing.T, group *model.Group, repo model.GroupRepository) {
	err := repo.Store(group)
	testset.FatalOnError(t, err, fmt.Sprintf("failed to store group %v", group))
}

func TestGroupStoreFind(t *testing.T) {
	repo := testGroupRepository(t)
	var expectedGroups []*model.Group

	for i := 1; i < 6; i++ {
		newG := model.BuildGroup(model.GroupID(i))
		expectedGroups = append(expectedGroups, newG)
		uploadTestGroup(t, newG, repo)
	}

	for _, group := range expectedGroups {
		_, err := repo.Find(group.ID)
		testset.FatalOnError(t, err, "failed to find group")

	}

	unexpectedGroupID := model.GroupID(42)
	if _, err := repo.Find(unexpectedGroupID); err == nil {
		t.Fatalf("Find should return err but got nil")
	}
}
