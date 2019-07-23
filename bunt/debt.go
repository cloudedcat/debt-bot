package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/tidwall/buntdb"
)

const indexDebt = "debts"
const prefixDebt = "debt::"

type debtRepository struct {
	db *buntdb.DB
}

// NewDebtRepository returns new instance of a BuntDB debt repository
func NewDebtRepository(dbName string) (model.DebtRepository, error) {
	db, err := buntdb.Open(dbName)
	if err != nil {
		return nil, err
	}

	return &debtRepository{db: db}, nil
}

func (r *debtRepository) FindAll(groupID model.GroupID) ([]*model.Debt, error) {
	_ = r.db.View(func(tx *buntdb.Tx) error {
		// tx.CreateIndex(debts, )
		tx.Ascend(indexDebt, func(key, value string) bool {
			return true
		})
		return nil
	})
	return nil, nil
}

func (r *debtRepository) Find(groupID model.GroupID, id int) (*model.Debt, error) {
	var val string
	err := r.db.View(func(tx *buntdb.Tx) error {
		var err error
		val, err = tx.Get(fmt.Sprintf("%s%d::%d", prefixDebt, groupID, id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return parseDebt(val)
}

func (r *debtRepository) Store(debt *model.Debt) error {
	wrpdDebt := (*debtWrapper)(debt)
	return r.db.Update(func(tx *buntdb.Tx) error {
		value, err := wrpdDebt.Value()
		if err != nil {
			return err
		}
		_, _, err = tx.Set(wrpdDebt.Key(), value, nil)
		return err
	})
}

func parseDebt(rawDebt string) (*model.Debt, error) {
	d := &model.Debt{}
	if err := json.Unmarshal([]byte(rawDebt), d); err != nil {
		return nil, err
	}
	return d, nil
}
