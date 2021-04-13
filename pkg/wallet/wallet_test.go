package wallet

import (
	"fmt"
	"testing"
)

func TestFindAccountByID_success(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	account, err = svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}

	account, err = svc.RegisterAccount("+992000000003")
	if err != nil {
		fmt.Println(err)
		return
	}
	// b := account.ID
	acc, err := svc.FindAccountByID(3)
	if err != nil {
		t.Errorf("%v: ", ErrAccountNotFound)
	}
	fmt.Println(account)
	fmt.Println(acc)
}
