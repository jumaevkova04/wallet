package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jumaevkova04/wallet/pkg/types"
)

func TestService_Reject(t *testing.T) {
	type args struct {
		paymentID string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Reject(tt.args.paymentID); (err != nil) != tt.wantErr {
				t.Errorf("Service.Reject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}
	return account, payments, nil
}

func TestService_Repeat_success(t *testing.T) {

	s := newTestService()

	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can't find payment by id, error = %v", err)
		return
	}

	if savedPayment.Status != types.PaymentStatusInProgress {
		t.Errorf("Repeat(): status didn't changed, payment = %v", savedPayment)
		return
	}
}

func TestService_FavoritePayment_success(t *testing.T) {

	s := newTestService()

	_, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]

	_, err = s.FindPaymentByID(favorite.ID)
	if err != nil {
		t.Errorf("FavoritePayment(): can't find favorite payment by id, error = %v", err)
		return
	}

	_, err = s.FavoritePayment(favorite.ID, "megafon")
	if err != nil {
		t.Errorf("FavoritePayment(): error = %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {

	// s := newTestService()
	// svc := &Service{}

	// _, favorites, err := s.addAccount(defaultTestAccount)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// favorite := favorites[0]
	// // _, err = s.FindPaymentByID(favorite.ID)
	// // if err != nil {
	// // 	t.Errorf("FavoritePayment(): can't find favorite payment by id, error = %v", err)
	// // 	return
	// // }

	// _, err = svc.FindFavoriteByID(favorite.ID)
	// if err != nil {
	// 	t.Errorf("FavoritePayment(): can't find favorite payment by id, error = %v", err)
	// 	return
	// }

	// _, err = svc.PayFromFavorite(favorite.ID)
	// if err != nil {
	// 	t.Errorf("FavoritePayment(): error = %v", err)
	// 	return
	// }

	s := newTestService()

	_, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]

	_, err = s.FindPaymentByID(favorite.ID)
	if err != nil {
		t.Errorf("FavoritePayment(): can't find favorite payment by id, error = %v", err)
		return
	}

	a, err := s.FavoritePayment(favorite.ID, "megafon")
	if err != nil {
		t.Errorf("FavoritePayment(): error = %v", err)
		return
	}

	_, err = s.FindFavoriteByID(a.ID)
	if err != nil {
		t.Errorf("FavoritePayment(): can't find favorite payment by id, error = %v", err)
		return
	}
	_, err = s.PayFromFavorite(a.ID)
	if err != nil {
		t.Errorf("FavoritePayment(): error = %v", err)
		return
	}

}
