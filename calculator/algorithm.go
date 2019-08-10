package calculator

import (
	"math"
	"sort"

	"github.com/cloudedcat/finance-bot/model"
)

const precision = 0.01

type pair struct {
	model.ParticipantID
	Amount float64
}

type calculatedDebt model.Debt

type pairs []pair

func (a pairs) Len() int      { return len(a) }
func (a pairs) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a pairs) Less(i, j int) bool {
	return less(a[i], a[j])
}

// use multikey sort to have always same order of pairs that get us always
// same result of debt calculation
func less(a, b pair) bool {
	if a.Amount != b.Amount {
		return a.Amount < b.Amount
	}
	return a.ParticipantID < b.ParticipantID
}

// calculate method decides whom and how much borrowers should return money.
// It has same complexity as sort.Sort O(n*log(n)). Probably that algorithm
// may get not optimal solutions and it's possible to reduce number of money
// transfers i.e. len of []calculatedDebt but it's sufficient for the bot purpose
func calculate(borrowers, lenders pairs) []calculatedDebt {
	sort.Sort(sort.Reverse(borrowers))
	sort.Sort(sort.Reverse(lenders))

	var debts []calculatedDebt
	var amount float64

	for len(borrowers) != 0 && len(lenders) != 0 {
		b, l := borrowers[0], lenders[0]
		amount = math.Min(b.Amount, l.Amount)
		amount = math.Round(amount/precision) * precision
		borrowers, lenders = shiftPair(borrowers), shiftPair(lenders)

		if math.Abs(b.Amount-l.Amount) < precision {
			// Do nothing if amounts are equal
		} else if b.Amount > l.Amount {
			// otherwise insert rest of the debt back into array
			b.Amount -= l.Amount
			borrowers = reverseInsertPair(borrowers, b)
		} else {
			l.Amount -= b.Amount
			lenders = reverseInsertPair(lenders, l)
		}
		d := calculatedDebt{
			BorrowerID: b.ParticipantID,
			LenderID:   l.ParticipantID,
			Amount:     amount,
		}
		debts = append(debts, d)
	}
	return debts
}

func reverseInsertPair(ps pairs, p pair) pairs {
	index := sort.Search(len(ps), func(i int) bool { return !less(p, ps[i]) })
	ps = append(ps, pair{})
	copy(ps[index+1:], ps[index:])
	ps[index] = p
	return ps
}

func shiftPair(ps pairs) pairs {
	return ps[1:]
}
