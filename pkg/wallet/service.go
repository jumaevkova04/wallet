package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jumaevkova04/wallet/pkg/types"
)

// ErrPhoneRegistered ..
var ErrPhoneRegistered = errors.New("phone already registered")

// ErrAmountMustBePositive ...
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")

// ErrAccountNotFound ...
var ErrAccountNotFound = errors.New("account not found")

// ErrNotEnoughBalance ...
var ErrNotEnoughBalance = errors.New("not enough balance")

// ErrPaymentNotFound ...
var ErrPaymentNotFound = errors.New("payment not found")

// ErrFavoriteNotFound ...
var ErrFavoriteNotFound = errors.New("favorite not found")

// Service ...
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

// RegisterAccount ...
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

// Deposit ...
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, cliceAccounts := range s.accounts {
		if cliceAccounts.ID == accountID {
			account = cliceAccounts
			break
		}
	}
	if account == nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж
	account.Balance += amount
	return nil
}

// Pay ...
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, cliceAccounts := range s.accounts {
		if cliceAccounts.ID == accountID {
			account = cliceAccounts
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

// FindAccountByID ...
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if accountID == account.ID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

// FindPaymentByID ...
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if paymentID == payment.ID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

// Reject ...
func (s *Service) Reject(paymentID string) error {

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

// Repeat ...
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, cliceAccounts := range s.accounts {
		if cliceAccounts.ID == payment.AccountID {
			account = cliceAccounts
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < payment.Amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= payment.Amount

	paymentID = uuid.New().String()
	payment = &types.Payment{
		ID:        paymentID,
		AccountID: payment.AccountID,
		Amount:    payment.Amount,
		Category:  payment.Category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

// FavoritePayment - создаёт избранное из конкретного платежа
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, cliceAccounts := range s.accounts {
		if cliceAccounts.ID == payment.AccountID {
			account = cliceAccounts
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < payment.Amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= payment.Amount
	paymentID = uuid.New().String()
	favorite := &types.Favorite{
		ID:        paymentID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

// FindFavoriteByID ...
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favoriteID == favorite.ID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}

// PayFromFavorite  - совершает платёж из конкретного избранного
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	if favorite.Amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, cliceAccounts := range s.accounts {
		if cliceAccounts.ID == favorite.AccountID {
			account = cliceAccounts
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < favorite.Amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= favorite.Amount

	favoriteID = uuid.New().String()
	payment := &types.Payment{
		ID:        favoriteID,
		AccountID: favorite.AccountID,
		Amount:    favorite.Amount,
		Category:  favorite.Category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

// Export ...
// type Export byte[]

// ExportToFile ...
func (s *Service) ExportToFile(path string) error {

	// service := &Service{}
	a := s.accounts

	cliceAccounts := []string{}
	var acc *types.Account
	var writeAccount string
	for i, account := range a {
		a[i] = account
		acc = a[i]

		id := strconv.Itoa(int(acc.ID))
		phone := acc.Phone
		balance := strconv.Itoa(int(acc.Balance))

		cliceAccounts = append(cliceAccounts, id, ";", string(phone), ";", balance, "|")

		writeAccount = strings.Join(cliceAccounts, "")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()

	_, err = f.Write([]byte(writeAccount))
	if err != nil {
		return err
	}
	return nil
}

// ImportFromFile ...
func (s *Service) ImportFromFile(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)
	// fmt.Print(data)

	accounts := strings.Split(data, "|")
	accounts = accounts[:len(accounts)-1]

	for _, account := range accounts {

		value := strings.Split(account, ";")
		id, err := strconv.Atoi(value[0])
		if err != nil {
			return err
		}
		phone := types.Phone(value[1])
		balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		editAccount := &types.Account{
			ID:      int64(id),
			Phone:   phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, editAccount)
		log.Print(account)
	}

	return nil
}
