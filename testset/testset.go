package testset

import (
	"testing"
	"time"

	"github.com/cloudedcat/debt-bot/model"
)

func FatalOnError(t *testing.T, err error, context string) {
	if err != nil {
		t.Fatalf("%s: %v", context, err)
	}
}

const Currency = "USD"

var GroupID model.GroupID = 42

func Debts() []*model.Debt {
	return []*model.Debt{
		{0, 15, "Bernard", Participants[0].ID, Participants[1].ID, time.Now(), Currency},
		{0, 25, "ChiniseCuisine", Participants[1].ID, Participants[2].ID, time.Now(), Currency},
		{0, 13, "BaoStore", Participants[3].ID, Participants[2].ID, time.Now(), Currency},
		{0, 13, "BaoStore", Participants[3].ID, Participants[1].ID, time.Now(), Currency},
		{0, 13, "BaoStore", Participants[2].ID, Participants[0].ID, time.Now(), Currency},
	}
}

var Participants = []*model.Participant{
	{12, "ck", "Louis", "C K"},
	{13, "noah", "Trevor", "Noah"},
	{14, "chapp", "Dave", "Chappelle"},
	{15, "murphy", "Eddie", "Murphy"},
	{16, "rock", "Chris", "Rock"},
}
