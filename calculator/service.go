package calculator

import (
	"github.com/cloudedcat/finance-bot/model"
)

type Service interface {
	AddDebt(groupID model.GroupID, debt model.Debt) error
	CalculateDebts(groupID model.GroupID) ([]FinalDebt, error)
}

// FinalDebt contains final sum that one participant owe another
// FinalDebt uses model.Debt as value object inside
type FinalDebt struct {
	Borrower *model.Participant
	Lender   *model.Participant
	model.Debt
}

type service struct {
	debts        model.DebtRepository
	participants model.ParticipantRepository
}

func NewService(debts model.DebtRepository, participants model.ParticipantRepository) Service {
	return &service{
		debts:        debts,
		participants: participants,
	}
}

func (s *service) AddDebt(groupID model.GroupID, debt model.Debt) error {
	var err error
	if debt.ID, err = s.debts.NextID(groupID); err != nil {
		return err
	}
	if err := debt.Validate(); err != nil {
		return err
	}
	return s.debts.Store(groupID, &debt)
}

// CalculateDebts is the main method of service and the main purpose of the application.
// It returns final debts for all participants that have debts, e.g.:
//   A owe to B 10 coins
//   B owe to C 15 coins
//   C owe to A 5 coins
// after calculation it returns:
//   A owe to B 5 coins
//   B owe to C 10 coins
// C owe nothing to A, so C is absent in FinalDebt array as borrower
func (s *service) CalculateDebts(groupID model.GroupID) ([]FinalDebt, error) {
	allDebts, err := s.debts.FindAll(groupID)
	if err != nil {
		return nil, err
	}
	partics, err := s.participants.FindAll(groupID)
	if err != nil {
		return nil, err
	}
	calculatedDebts := s.calculate(allDebts)
	return s.composeFinalDebts(calculatedDebts, partics), nil
}

func (s *service) calculate(debts []*model.Debt) []calculatedDebt {
	debtSum := make(map[model.ParticipantID]float64)
	for _, debt := range debts {
		debtSum[debt.BorrowerID] += debt.Amount
		debtSum[debt.LenderID] -= debt.Amount
	}
	var borrorwers, lenders []pair
	for id, amount := range debtSum {
		if amount > 0 {
			borrorwers = append(borrorwers, pair{id, amount})
		} else if amount < 0 {
			amount = -amount
			lenders = append(lenders, pair{id, amount})
		}
	}
	return calculate(borrorwers, lenders)
}

func (s *service) composeFinalDebts(debts []calculatedDebt, partics model.Participants) (final []FinalDebt) {
	particMap := partics.AsMap()
	for _, debt := range debts {
		fDebt := FinalDebt{
			Borrower: particMap[debt.BorrowerID],
			Lender:   particMap[debt.LenderID],
			Debt:     model.Debt(debt),
		}
		final = append(final, fDebt)
	}
	return
}
