package bunt

import (
	"encoding/json"
	"fmt"

	"github.com/cloudedcat/finance-bot/model"
)

type debtWrapper model.Debt

func (d debtWrapper) Key() string {
	return fmt.Sprintf("%s%d::%d", prefixDebt, d.GroupID, d.ID)
}

func (d debtWrapper) Index() string {
	return fmt.Sprintf("%s%d", prefixDebt, d.GroupID)
}

// func (d debtWrapper) Pattern() string {
// 	return fmt.Sprintf("%s%d", prefixDebt, d.GroupID)
// }

func (d debtWrapper) Value() (string, error) {
	jsonb, err := json.Marshal(d)
	return string(jsonb), err
}
