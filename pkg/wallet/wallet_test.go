package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jumaevkova04/wallet/pkg/types"

	"github.com/google/uuid"
)

func TestFindAccountByID_success(t *testing.T) {
	svc := &Service{}

	_, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = svc.RegisterAccount("+992000000002")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = svc.RegisterAccount("+992000000003")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = svc.FindAccountByID(3)
	if err != nil {
		t.Errorf("%v: ", ErrAccountNotFound)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
		return
	}

	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}

}

func TestService_FindPaymentByID_fail(t *testing.T) {

	s := newTestService()

	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}

func TestReject_success(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}

	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}

	if savedAccount.Balance != (defaultTestAccount.balance) {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}
