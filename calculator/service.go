//go:generate mockgen -destination mock_calculator/mock_calculator.go github.com/cloudedcat/debt-bot/calculator Service

package calculator

import (
	"errors"
	"time"

	"github.com/cloudedcat/debt-bot/model"
)

type Service interface {
	AddDebtsByAliases(groupID model.GroupID, debts ...DebtWithAliases) error
	// AddDebts(groupID model.GroupID, debts ...*model.Debt) error
	CalculateDebts(groupID model.GroupID) ([]DetailedDebt, error)
	FindDebts(groupID model.GroupID, pID model.ParticipantID) ([]DetailedDebt, error)
	ClearDebts(groupID model.GroupID) error
}

var ErrParticipantNotFound = errors.New("participant not found")

// DetailedDebt contains Borrower and Lender entities inside
type DetailedDebt struct {
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

func (s *service) AddDebts(groupID model.GroupID, debts ...*model.Debt) error {
	var err error
	for _, debt := range debts {
		if debt.ID, err = s.debts.NextID(groupID); err != nil {
			return err
		}
		if err = debt.Validate(); err != nil {
			return err
		}
	}
	return s.debts.Store(groupID, debts...)
}

type DebtWithAliases struct {
	Amount   float64
	Tag      string
	Borrower model.Alias
	Lender   model.Alias
}

func (s *service) AddDebtsByAliases(groupID model.GroupID, aliasDebts ...DebtWithAliases) error {
	partics, err := s.participants.FindAll(groupID)
	if err != nil {
		return err
	}
	aliasMap := partics.AsAliasMap()
	timeNow := time.Now()

	var debts []*model.Debt
	for _, aliasDebt := range aliasDebts {
		borrower, lender := aliasMap[aliasDebt.Borrower], aliasMap[aliasDebt.Lender]
		if borrower == nil || lender == nil {
			return ErrParticipantNotFound
		}
		debt := &model.Debt{
			LenderID:   lender.ID,
			BorrowerID: borrower.ID,
			Tag:        aliasDebt.Tag,
			Amount:     aliasDebt.Amount,
			Date:       timeNow,
		}
		if debt.ID, err = s.debts.NextID(groupID); err != nil {
			return err
		}
		if err = debt.Validate(); err != nil {
			return err
		}
		debts = append(debts, debt)
	}
	return s.debts.Store(groupID, debts...)
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
func (s *service) CalculateDebts(groupID model.GroupID) ([]DetailedDebt, error) {
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

func (s *service) composeFinalDebts(debts []calculatedDebt, partics model.Participants) (final []DetailedDebt) {
	particMap := partics.AsMap()
	for _, debt := range debts {
		fDebt := DetailedDebt{
			Borrower: particMap[debt.BorrowerID],
			Lender:   particMap[debt.LenderID],
			Debt:     model.Debt(debt),
		}
		final = append(final, fDebt)
	}
	return
}

func (s *service) FindDebts(groupID model.GroupID, pID model.ParticipantID) ([]DetailedDebt, error) {
	allDebts, err := s.debts.FindAll(groupID)
	if err != nil {
		return nil, err
	}

	partics, err := s.participants.FindAll(groupID)
	if err != nil {
		return nil, err
	}
	particMap := partics.AsMap()
	selected := []DetailedDebt{}
	for _, d := range allDebts {
		if d.BorrowerID == pID || d.LenderID == pID {
			detailedDebt := DetailedDebt{
				Borrower: particMap[d.BorrowerID],
				Lender:   particMap[d.LenderID],
				Debt:     *d,
			}
			selected = append(selected, detailedDebt)
		}
	}
	return selected, nil
}

func (s *service) ClearDebts(groupID model.GroupID) error {
	return s.debts.Clear(groupID)
}
