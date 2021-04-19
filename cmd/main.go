package main

import (
	"github.com/jumaevkova04/wallet/pkg/wallet"
)

func main() {
	s := &wallet.Service{}

	s.RegisterAccount("+992000000001")
	s.RegisterAccount("+992000000002")
	s.RegisterAccount("+992000000003")
	s.ExportToFile("dump.txt")

}
