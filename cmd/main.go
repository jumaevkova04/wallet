package main

import (
	"fmt"

	"github.com/jumaevkova04/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(account.ID)

	account, err = svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(account.ID)

	// err = svc.Deposit(account.ID, 10)
	// if err != nil {
	// 	switch err {
	// 	case wallet.ErrAmountMustBePositive:
	// 		fmt.Println("Сумма должно быть положительной")
	// 	case wallet.ErrAccountNotFound:
	// 		fmt.Println("Аккаунт пользователя не найден")
	// 	}
	// 	return
	// }
	// fmt.Println(account.Balance)

	acc, err := svc.FindAccountByID(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(acc.ID)

}
