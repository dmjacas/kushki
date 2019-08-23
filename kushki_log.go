package kushki

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

// KushkiRequestLog Model
type KushkiRequestLog struct {
	ID             int       `json:"id" gorm:"PRIMARY_KEY; AUTO_INCREMENT;size:11" `
	Active         bool      `json:"active"`
	name           string    `json:"reference"`
	number         string    `json:"allResponse" gorm:"size:5550"`
	expiryMonth    string    `json:"expiration" gorm:"size:2550"`
	expiryYear     string    `json:"ipadres"`
	cvv            string    `json:"returnUrl" gorm:"size:550"`
	CancelURL      string    `json:"cancelUrl" gorm:"size:550"`
	SkipResult     bool      `json:"skipResult" `
	NoBuyerFill    bool      `json:"noBuyerFill"`
	CaptureAddress bool      `json:"captureAddress"`
	PaymentMethod  bool      `json:"paymentMethod"`
	Fields         string    `json:"fields" gorm:"size:2550"`
	RequestID      string    `json:"requestId"`
	ProcessURL     string    `json:"processUrl" gorm:"size:250"`
	Message        string    `json:"message" gorm:"size:250"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	DeletedAt      null.Time `json:"-"`
}
