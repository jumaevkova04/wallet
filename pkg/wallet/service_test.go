package wallet

import (
	"fmt"
	"testing"

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
