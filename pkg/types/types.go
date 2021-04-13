package types

// Money предстовляет собой денежную сумму в минимальных еденицах (центы, копейки, дирамы и е.д.).
type Money int64

// PaymentCategory  предстовляет собой категорию, в которой был совершён платёж (авто, аптеки, рестораны и т.д.).
type PaymentCategory string

// PaymentStatus представляет собой статус платежа.
type PaymentStatus string

// Предопределённые статусы платежей.
const (
	PaymentStatusOk         PaymentStatus = "OK"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

// Payment предстовляет информацию о платеже.
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

// Phone ...
type Phone string

// Account представляет информацию о счёте пользователья.
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
