package postgre

import (
	"github.com/lib/pq"
)

//StringArray postgres array type support for xorm
type StringArray []string

func (a *StringArray) FromDB(bts []byte) error {
	pqArray := new(pq.StringArray)
	err := pqArray.Scan(bts)
	if err != nil {
		return err
	}
	*a = []string(*pqArray)
	return nil
}

func (a StringArray) ToDB() ([]byte, error) {
	pqArray := pq.StringArray(a)
	v, err := pqArray.Value()
	return []byte(v.(string)), err
}
