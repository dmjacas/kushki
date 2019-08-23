package kushki

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

// KushkiRequestLog Model
type KushkiRequestLog struct {
	ID                        int       `json:"id" gorm:"PRIMARY_KEY; AUTO_INCREMENT;size:11" `
	Proces                    string    `json:"proces" gorm:"size:250"`
	Code                      string    `json:"code" gorm:"size:2240"`
	Message                   string    `json:"message" gorm:"size:2240"`
	AcquirerBank              string    `json:"acquirerBank" gorm:"size:2240"`
	ApprovedTransactionAmount float64   `json:"approvedTransactionAmount"`
	SubtotalIva               float64   `json:"subtotalIva"`
	SubtotalIva0              float64   `json:"subtotalIva0"`
	ResponseCode              string    `json:"responseCode" gorm:"size:250"`
	TransactionType           string    `json:"transactionType" gorm:"size:250"`
	ApprovalCode              string    `json:"approvalCode" gorm:"size:250"`
	TransactionStatus         string    `json:"transactionStatus" gorm:"size:250"`
	CurrencyCode              string    `json:"currencyCode" gorm:"size:250"`
	ProcessorID               string    `json:"processorId" gorm:"size:250"`
	TransactionID             string    `json:"transactionId" gorm:"size:250"`
	ResponseText              string    `json:"responseText" gorm:"size:250"`
	CardHolderName            string    `json:"cardHolderName" gorm:"size:250"`
	LastFourDigits            string    `json:"lastFourDigits" gorm:"size:250"`
	BinCard                   string    `json:"binCard" gorm:"size:250"`
	PaymentBrand              string    `json:"paymentBrand" gorm:"size:250"`
	RequestAmount             float64   `json:"requestAmount"`
	IvaValue                  float64   `json:"ivaValue"`
	MerchantName              string    `json:"merchantName" gorm:"size:250"`
	ProcessorName             string    `json:"processorName" gorm:"size:250"`
	TransactionReference      string    `json:"transactionReference" gorm:"size:250"`
	CreatedAt                 time.Time `json:"createdAt"`
	UpdatedAt                 time.Time `json:"updatedAt"`
	DeletedAt                 null.Time `json:"-"`
}
