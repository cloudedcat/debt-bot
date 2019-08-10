package calculator

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func getTestPairs() pairs {
	return pairs{
		{1, 10},
		{2, 7},
		{1, 7},
		{5, 8},
	}
}

func TestSort(t *testing.T) {
	testPairs := getTestPairs()
	sort.Sort(testPairs)
	expectedPair, index := pair{1, 7}, 0
	if testPairs[index] != expectedPair {
		t.Fatalf("unexpected pair on index %d: expected %v, but got %v", index, expectedPair, testPairs[index])
	}
}

func TestReverseInsert(t *testing.T) {
	testPairs := getTestPairs()
	sort.Sort(sort.Reverse(testPairs))
	testPairs = reverseInsertPair(testPairs, pair{3, 7})
	expectedPair, expectedlength := pair{3, 7}, 5
	index := 2
	if expectedlength != len(testPairs) {
		t.Fatalf("unexpected size of testPairs: expected %d, got %d", expectedlength, len(testPairs))
	}
	if testPairs[index] != expectedPair {
		t.Fatalf("unexpected pair on index %d: expected %v, got %v", index, expectedPair, testPairs[index])
	}
}

func TestShiftPair(t *testing.T) {
	testPairs := getTestPairs()
	sort.Sort(sort.Reverse(testPairs))
	testPairs = shiftPair(testPairs)
	expectedPair, expectedlength := pair{5, 8}, 3
	index := 0

	if expectedlength != len(testPairs) {
		t.Fatalf("unexpected size of testPairs: expected %d, got %d", expectedlength, len(testPairs))
	}
	if testPairs[index] != expectedPair {
		t.Fatalf("unexpected pair on index %d: expected %v, got %v", index, expectedPair, testPairs[index])
	}
}

var calculateTestCases = []struct {
	borrowers pairs
	lenders   pairs
	result    []calculatedDebt
}{
	{
		borrowers: pairs{{1, 1}, {2, 2}, {3, 5}, {4, 5.5}},
		lenders:   pairs{{5, 3}, {6, 4}, {7, 6}, {8, 0.5}},
		result: []calculatedDebt{
			{BorrowerID: 4, LenderID: 7, Amount: 5.5},
			{BorrowerID: 3, LenderID: 6, Amount: 4},
			{BorrowerID: 2, LenderID: 5, Amount: 2},
			{BorrowerID: 3, LenderID: 5, Amount: 1},
			{BorrowerID: 1, LenderID: 8, Amount: 0.5},
			{BorrowerID: 1, LenderID: 7, Amount: 0.5},
		},
	},
	{
		borrowers: pairs{{1, 2}, {2, 3}, {3, 5}},
		lenders:   pairs{{4, 3}, {5, 3}, {6, 4}},
		result: []calculatedDebt{
			{BorrowerID: 3, LenderID: 6, Amount: 4},
			{BorrowerID: 2, LenderID: 5, Amount: 3},
			{BorrowerID: 1, LenderID: 4, Amount: 2},
			{BorrowerID: 3, LenderID: 4, Amount: 1},
		},
	},
}

func TestCalculateAlgorithm(t *testing.T) {
	for _, testCase := range calculateTestCases {
		result := calculate(testCase.borrowers, testCase.lenders)
		testPrintDebts(t, result)
		if diff := cmp.Diff(testCase.result, result); diff != "" {
			t.Fatalf("Wrong result (-expected, +got):\n%s", diff)
		}
	}
}

func testPrintDebts(t *testing.T, ds []calculatedDebt) {
	for _, d := range ds {
		testPrintDebt(t, d)
	}
}
func testPrintDebt(t *testing.T, d calculatedDebt) {
	t.Logf("%v -> %v - %f coins ", d.BorrowerID, d.LenderID, d.Amount)
}
