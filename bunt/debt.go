package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/tidwall/buntdb"
)

func prefixDebt(groupID model.GroupID) string {
	return fmt.Sprintf("debt%s%d", sep, int(groupID))
}

func indexDebt(groupID model.GroupID) string {
	return prefixDebt(groupID)
}

type debtRepository struct {
	db *buntdb.DB
}

// NewDebtRepository returns new instance of a BuntDB debt repository
func NewDebtRepository(db *buntdb.DB) model.DebtRepository {
	return &debtRepository{db: db}
}

func (r *debtRepository) FindAll(groupID model.GroupID) ([]*model.Debt, error) {
	var debts []*model.Debt
	err := r.db.View(func(tx *buntdb.Tx) error {
		var pErr error
		txErr := tx.Ascend(indexDebt(groupID), func(_, raw string) bool {
			var parsed *model.Debt
			if parsed, pErr = r.parse(raw); pErr != nil {
				return false
			}
			debts = append(debts, parsed)
			return true
		})
		if pErr != nil {
			return pErr
		}
		return txErr
	})
	if err != nil {
		return nil, err
	}
	return debts, nil
}

func (r *debtRepository) Find(groupID model.GroupID, id model.DebtID) (*model.Debt, error) {
	var raw string
	err := r.db.View(func(tx *buntdb.Tx) error {
		var err error
		raw, err = tx.Get(r.key(groupID, id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return r.parse(raw)
}

func (r *debtRepository) Store(groupID model.GroupID, debt *model.Debt) error {
	return r.db.Update(func(tx *buntdb.Tx) error {
		raw, err := r.compose(debt)
		if err != nil {
			return err
		}
		_, _, err = tx.Set(r.key(groupID, debt.ID), raw, nil)
		return err
	})
}

func (r *debtRepository) key(groupID model.GroupID, id model.DebtID) string {
	return fmt.Sprintf("%s%s%d", prefixDebt(groupID), sep, int(id))
}

func (r *debtRepository) NextID(groupID model.GroupID) (model.DebtID, error) {
	// NYI
	return 0, nil
}

func (r *debtRepository) parse(rawDebt string) (*model.Debt, error) {
	debt := &model.Debt{}
	if err := json.Unmarshal([]byte(rawDebt), debt); err != nil {
		return nil, err
	}
	return debt, nil
}

func (r *debtRepository) compose(debt *model.Debt) (string, error) {
	b, err := json.Marshal(debt)
	return string(b), err
}
