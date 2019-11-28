package kushki

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// KushkiURL kushki url payment
var KushkiURL string

// KushkiPublicKey public key
var KushkiPublicKey string

//KushkiPrivateKey private key
var KushkiPrivateKey string

// P2PDB db object
var P2PDB *gorm.DB
var db *dBConfig

// Config configure payment library
// KushkiURL Kushki url
// KushkiPublicKey Kushki public key
// KushkiPrivateKey Kushki private key

// dbCharset db Charset
// dbDialect db Dialect
// dbName dn name
// dbUsername db username
// dbPassword db password
// Expiration time in minutes that the payment request lasts

// Config pay db connectiong
func Config(urlKushkiURLPayment, ParamKushkiPublicKey, ParamKushkiPrivateKey, dbCharset, dbDialect, dbName, dbUsername, dbHost, dbPort, dbPassword string, Expiration int) {
	KushkiURL = "https://api-uat.kushkipagos.com/card/v1/" //urlKushkiURLPayment
	KushkiPublicKey = "20000000107193962000"               //ParamKushkiPublicKey
	KushkiPrivateKey = "20000000102569300000"              //ParamKushkiPrivateKey
	db := &dBConfig{
		Dialect:  dbDialect,
		Username: dbUsername,
		Password: dbPassword,
		Host:     dbHost,
		Name:     dbName,
		Port:     dbPort,
		Charset:  dbCharset,
	}
	P2PDB, _ = Connect(db)
	migration()

}

// migration  create table if not exist
func migration() {
	pingErr := P2PDB.DB().Ping()
	if pingErr != nil {
		fmt.Println(pingErr)
	} else {
		P2PDB.AutoMigrate(&KushkiRequestLog{})
	}
}

// dBConfig database config structure
type dBConfig struct {
	Host     string `default:"localhost"`
	Port     string `default:"3306"`
	Dialect  string `default:"mysql"`
	Username string
	Password string
	Name     string
	Charset  string `default:"utf8"`
}

// Connect handles the connection to the database and returns it
func Connect(config *dBConfig) (*gorm.DB, error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		config.Charset)
	db, err := gorm.Open(config.Dialect, dbURI)
	if err != nil {
		log.Fatalln("db connect", err)
	}

	return db, nil
}

// RequestTokenCard request token card
func RequestTokenCard(par *Request) (*KushkiResponse, error) {

	jsonRequest, err := json.Marshal(par)
	if err != nil {
		return nil, errors.New("error to JSON encode the body request")
	}
	client := http.Client{}
	request, err := http.NewRequest("POST", KushkiURL+"tokens", bytes.NewBuffer([]byte(jsonRequest)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("Public-Merchant-Id", KushkiPublicKey)

	response, err := client.Do(request)

	if err != nil {
		return nil, errors.New("error in the http call")
	}
	dataResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	res := &KushkiResponse{}
	errp := json.Unmarshal(dataResp, res)
	if errp != nil {
		fmt.Println(err)
	}
	return res, nil
}

// RequestCharges add charde
func RequestCharges(par *ChargeParams) (*ChargeResponse, error) {
	jsonRequest, err := json.Marshal(par)
	if err != nil {
		return nil, errors.New("error to JSON encode the body request")
	}
	client := http.Client{}
	request, err := http.NewRequest("POST", KushkiURL+"charges", bytes.NewBuffer([]byte(jsonRequest)))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Set("Private-Merchant-Id", KushkiPrivateKey)

	response, err := client.Do(request)
	if err != nil {
		return nil, errors.New("error in the http call")
	}
	dataResp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	res := &ChargeResponse{}
	errp := json.Unmarshal(dataResp, res)
	if errp != nil {
		fmt.Println(err)
	}
	return res, nil
}

// CancelTransaction cancel transaction
func CancelTransaction(ticketNumber string, par *PreAuthorizationParams) (*KushkiResponse, error) {
	jsonRequest, err := json.Marshal(par)
	if err != nil {
		return nil, errors.New("error to JSON encode the body request")
	}
	req, err := http.NewRequest("DELETE", KushkiURL+"charges/"+ticketNumber, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Merchant-Id", KushkiPrivateKey)
	if err != nil {
		return nil, errors.New("error in the http call")
	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error in the http call")
	}

	dataResp, err := ioutil.ReadAll(response.Body)
	var KushKiResponse KushkiResponse
	if err = json.Unmarshal(dataResp, &KushKiResponse); err != nil {
		return nil, errors.New("error in the return values of the http call")
	}

	return &KushKiResponse, nil
}

// ReimburseTransaction reimburse transaction
func ReimburseTransaction(ticketNumber string, par *PreAuthorizationParams) (*KushkiResponse, error) {
	jsonRequest, err := json.Marshal(par)
	if err != nil {
		return nil, errors.New("error to JSON encode the body request")
	}
	req, err := http.NewRequest("DELETE", KushkiURL+"refund/"+ticketNumber, bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Merchant-Id", KushkiPrivateKey)
	if err != nil {
		return nil, errors.New("error in the http call")
	}
	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error in the http call")
	}

	dataResp, err := ioutil.ReadAll(response.Body)
	var KushKiResponse KushkiResponse
	if err = json.Unmarshal(dataResp, &KushKiResponse); err != nil {
		return nil, errors.New("error in the return values of the http call")
	}

	return &KushKiResponse, nil
}

//PreAuthorizationPayment preauthorization payment
func PreAuthorizationPayment(par *PreAuthorizationParams) (*KushkiResponse, error) {
	jsonRequest, err := json.Marshal(par)
	if err != nil {
		return nil, errors.New("error to JSON encode the body request")
	}

	req, err := http.NewRequest("POST", KushkiURL+"preAuthorization", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Merchant-Id", KushkiPrivateKey)
	if err != nil {
		return nil, errors.New("error in the http call")
	}
	client := &http.Client{Timeout: time.Second * 100}
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error in the http call")
	}

	dataResp, err := ioutil.ReadAll(response.Body)
	var KushKiResponse KushkiResponse
	if err = json.Unmarshal(dataResp, &KushKiResponse); err != nil {
		return nil, errors.New("error in the return values of the http call")
	}

	kushkiRequestLog := &KushkiRequestLog{
		Proces: "preAuthorization",
		Code:   KushKiResponse.Code,
	}
	tx := P2PDB.Begin()
	if result := tx.Create(&kushkiRequestLog); result.Error != nil {
		tx.Rollback()
		return nil, errors.New("error in saving the data")
	}
	if result := tx.Commit(); result.Error != nil {
		tx.Rollback()
		return nil, errors.New("error in saving the data")
	}
	return &KushKiResponse, nil
}

// CaptureAuthorizationPayment capture payment
func CaptureAuthorizationPayment(par *CaptureParams) (*KushkiResponse, error) {

	jsonRequest, err := json.Marshal(par)
	if err != nil {
		return nil, errors.New("error to JSON encode the body request")
	}
	req, err := http.NewRequest("POST", KushkiURL+"capture", bytes.NewBuffer(jsonRequest))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Private-Merchant-Id", KushkiPrivateKey)
	if err != nil {
		return nil, errors.New("error in the http call")
	}
	client := &http.Client{Timeout: time.Second * 10}
	fmt.Println(KushkiPrivateKey)
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error in the http call")
	}

	dataResp, err := ioutil.ReadAll(response.Body)
	var KushKiResponse KushkiResponse
	if err = json.Unmarshal(dataResp, &KushKiResponse); err != nil {
		return nil, errors.New("error in the return values of the http call")
	}

	return &KushKiResponse, nil
}

// Card card structure
type Card struct {
	Name        string `json:"name"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiryMonth"`
	ExpiryYear  string `json:"expiryYear"`
	CVV         string `json:"cvv"`
}

// Request request structure
type Request struct {
	Card        *Card   `json:"card"`
	TotalAmount float64 `json:"totalAmount"`
	Currency    string  `json:"currency"`
}

// Amount amount structure
type Amount struct {
	SubtotalIVA  float64 `json:"subtotalIva"`
	SubtotalIva0 float64 `json:"subtotalIva0"`
	Ice          float64 `json:"ice"`
	Iva          float64 `json:"iva"`
	Currency     string  `json:"currency"`
}

//Deferred struct
type Deferred struct {
	GraceMonths string `json:"graceMonths,omitempty"`
	CreditType  string `json:"creditType,omitempty"`
	Months      int    `json:"months,omitempty"`
}

//Details response detail
type Details struct {
	ApprovalCode              string      `json:"approvalCode,omitempty"`
	ApprovedTransactionAmount float64     `json:"approvedTransactionAmount,omitempty"`
	BinCard                   string      `json:"binCard,omitempty"`
	CardHolderName            string      `json:"cardHolderName,omitempty"`
	CardType                  string      `json:"cardType,omitempty"`
	LastFourDigitsOfCard      string      `json:"lastFourDigitsOfCard,omitempty"`
	MerchantName              string      `json:"merchantName,omitempty"`
	ProcessorName             string      `json:"processorName,omitempty"`
	Recap                     string      `json:"recap,omitempty"`
	ResponseCode              string      `json:"responseCode,omitempty"`
	ResponseText              string      `json:"responseText,omitempty"`
	TicketNumber              string      `json:"ticketNumber,omitempty"`
	TransactionID             string      `json:"transactionId,omitempty"`
	Token                     string      `json:"token,omitempty"`
	FullResponse              bool        `json:"fullResponse,omitempty" `
	AcquirerBank              string      `json:"acquirerBank,omitempty"`
	IP                        string      `json:"ip,omitempty"`
	MaskedCardNumber          string      `json:"maskedCardNumber,omitempty"`
	SubtotalIva               string      `json:"subtotalIva,omitempty" `
	SubtotalIva0              string      `json:"subtotalIva0,omitempty"`
	Created                   string      `json:"created,omitempty"`
	TransactionType           string      `json:"transactionType,omitempty"`
	TransactionStatus         string      `json:"transactionStatus,omitempty"`
	SyncMode                  string      `json:"syncMode,omitempty"`
	CurrencyCode              string      `json:"currencyCode,omitempty"`
	MerchantID                string      `json:"merchantId,omitempty"`
	ProcessorID               string      `json:"processorId,omitempty"`
	LastFourDigits            string      `json:"lastFourDigits,omitempty"`
	PaymentBrand              string      `json:"paymentBrand,omitempty"`
	IceValue                  string      `json:"iceValue,omitempty"`
	RequestAmount             string      `json:"requestAmount,omitempty"`
	IvaValue                  string      `json:"ivaValue,omitempty"`
	ProcessorBankName         string      `json:"processorBankName,omitempty"`
	TransactionReference      string      `json:"transactionReference,omitempty"`
	Errors                    interface{} `json:"errors,omitempty"`
}

//BinInfo modal
type BinInfo struct {
	Bank string `json:"bank,omitempty"`
	Type string `json:"type,omitempty"`
}

// Metadata  params
type Metadata struct {
	CustomerID string `json:"customerId"`
	ContractID string `json:"contractID"`
}

//PreAuthorizationParams paramas
type PreAuthorizationParams struct {
	Token        string    `json:"token"`
	Amount       *Amount   `json:"amount"`
	Metadata     *Metadata `json:"metadata"`
	FullResponse bool      `json:"fullResponse"`
}

// KushkiResponse response
type KushkiResponse struct {
	Code         string   `json:"code"`
	Details      *Details `json:"details"`
	Message      string   `json:"message"`
	BinInfo      *BinInfo `json:"binInfo,omitempty"`
	TicketNumber string   `json:"ticketNumber,omitempty"`
	Token        string   `json:"token,omitempty"`
}

// CaptureParams struct
type CaptureParams struct {
	TicketNumber string      `json:"ticketNumber,omitempty"`
	Amount       *Amount     `json:"amount,omitempty"`
	FullResponse interface{} `json:"fullResponse"`
	Metadata     bool        `json:"matadata,omitempty"`
}

// ChargeParams struct
type ChargeParams struct {
	Token        string    `json:"token,omitempty"`
	Amount       *Amount   `json:"amount,omitempty"`
	Deferred     *Deferred `json:"deferred,omitempty"`
	FullResponse bool      `json:"fullResponse"`
	Metadata     *Metadata `json:"matadata,omitempty"`
}

// Token struct
type Token struct {
	Token string `json:"token,omitempty"`
}

//ChargeResponse struct
type ChargeResponse struct {
	TicketNumber string `json:"ticketNumber"`
	Details      struct {
		Token                     string `json:"token"`
		FullResponse              bool   `json:"fullResponse"`
		Recap                     string `json:"recap"`
		AcquirerBank              string `json:"acquirerBank"`
		IP                        string `json:"ip"`
		MaskedCardNumber          string `json:"maskedCardNumber"`
		ApprovedTransactionAmount int    `json:"approvedTransactionAmount"`
		SubtotalIva               int    `json:"subtotalIva"`
		SubtotalIva0              int    `json:"subtotalIva0"`
		Created                   int64  `json:"created"`
		ResponseCode              string `json:"responseCode"`
		TransactionType           string `json:"transactionType"`
		ApprovalCode              string `json:"approvalCode"`
		TransactionStatus         string `json:"transactionStatus"`
		SyncMode                  string `json:"syncMode"`
		CurrencyCode              string `json:"currencyCode"`
		MerchantID                string `json:"merchantId"`
		ProcessorID               string `json:"processorId"`
		TransactionID             string `json:"transactionId"`
		ResponseText              string `json:"responseText"`
		CardHolderName            string `json:"cardHolderName"`
		LastFourDigits            string `json:"lastFourDigits"`
		BinCard                   string `json:"binCard"`
		PaymentBrand              string `json:"paymentBrand"`
		IceValue                  int    `json:"iceValue"`
		RequestAmount             int    `json:"requestAmount"`
		IvaValue                  int    `json:"ivaValue"`
		MerchantName              string `json:"merchantName"`
		ProcessorName             string `json:"processorName"`
		ProcessorBankName         string `json:"processorBankName"`
		TransactionReference      string `json:"transactionReference"`
		BinInfo                   struct {
			Bank string `json:"bank"`
			Type string `json:"type"`
		} `json:"binInfo"`
	} `json:"details"`
}
