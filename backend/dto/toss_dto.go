package dto

import "time"

type Payment struct {
	Version             string               `json:"version"`
	PaymentKey          string               `json:"paymentKey"`
	Type                string               `json:"type"`
	OrderId             string               `json:"orderId"`
	OrderName           string               `json:"orderName"`
	MId                 string               `json:"mId"`
	Currency            string               `json:"currency"`
	Method              string               `json:"method"`
	TotalAmount         float64              `json:"totalAmount"`
	BalanceAmount       float64              `json:"balanceAmount"`
	Status              string               `json:"status"`
	RequestedAt         string               `json:"requestedAt"`
	ApprovedAt          string               `json:"approvedAt"`
	UseEscrow           bool                 `json:"useEscrow"`
	LastTransactionKey  *string              `json:"lastTransactionKey,omitempty"`
	SuppliedAmount      float64              `json:"suppliedAmount"`
	Vat                 float64              `json:"vat"`
	CultureExpense      bool                 `json:"cultureExpense"`
	TaxFreeAmount       float64              `json:"taxFreeAmount"`
	TaxExemptionAmount  int                  `json:"taxExemptionAmount"`
	Cancels             []Cancel             `json:"cancels,omitempty"`
	IsPartialCancelable bool                 `json:"isPartialCancelable"`
	Card                *CardInfo            `json:"card,omitempty"`
	VirtualAccount      *VirtualAccountInfo  `json:"virtualAccount,omitempty"`
	MobilePhone         *MobilePhoneInfo     `json:"mobilePhone,omitempty"`
	GiftCertificate     *GiftCertificateInfo `json:"giftCertificate,omitempty"`
	Transfer            *TransferInfo        `json:"transfer,omitempty"`
	Receipt             *ReceiptInfo         `json:"receipt,omitempty"`
	Checkout            *CheckoutInfo        `json:"checkout,omitempty"`
	EasyPay             *EasyPayInfo         `json:"easyPay,omitempty"`
	Country             string               `json:"country,omitempty"`
	Failure             *FailureInfo         `json:"failure,omitempty"`
	CashReceipt         *CashReceiptInfo     `json:"cashReceipt,omitempty"`
	CashReceipts        []CashReceiptInfo    `json:"cashReceipts,omitempty"`
	Discount            *DiscountInfo        `json:"discount,omitempty"`
}

type CardInfo struct {
	Amount                float64 `json:"amount"`
	IssuerCode            string  `json:"issuerCode"`
	AcquirerCode          *string `json:"acquirerCode,omitempty"`
	Number                string  `json:"number"`
	InstallmentPlanMonths int     `json:"installmentPlanMonths"`
	ApproveNo             string  `json:"approveNo"`
	UseCardPoint          bool    `json:"useCardPoint"`
	CardType              string  `json:"cardType"`
	OwnerType             string  `json:"ownerType"`
	AcquireStatus         string  `json:"acquireStatus"`
	IsInterestFree        bool    `json:"isInterestFree"`
	InterestPayer         *string `json:"interestPayer,omitempty"`
}

type VirtualAccountInfo struct {
	AccountType          string                    `json:"accountType"`
	AccountNumber        string                    `json:"accountNumber"`
	BankCode             string                    `json:"bankCode"`
	CustomerName         string                    `json:"customerName"`
	DueDate              string                    `json:"dueDate"`
	RefundStatus         string                    `json:"refundStatus"`
	Expired              bool                      `json:"expired"`
	SettlementStatus     string                    `json:"settlementStatus"`
	RefundReceiveAccount *RefundReceiveAccountInfo `json:"refundReceiveAccount,omitempty"`
	Secret               *string                   `json:"secret,omitempty"`
}

type MobilePhoneInfo struct {
	CustomerMobilePhone string `json:"customerMobilePhone"`
	SettlementStatus    string `json:"settlementStatus"`
	ReceiptUrl          string `json:"receiptUrl"`
}

type GiftCertificateInfo struct {
	ApproveNo        string `json:"approveNo"`
	SettlementStatus string `json:"settlementStatus"`
}

type TransferInfo struct {
	BankCode         string `json:"bankCode"`
	SettlementStatus string `json:"settlementStatus"`
}

type ReceiptInfo struct {
	Url string `json:"url"`
}

type CheckoutInfo struct {
	Url string `json:"url"`
}

type EasyPayInfo struct {
	Provider       string  `json:"provider"`
	Amount         float64 `json:"amount"`
	DiscountAmount float64 `json:"discountAmount"`
}

type FailureInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CashReceiptInfo struct {
	Type                   string       `json:"type"`
	ReceiptKey             string       `json:"receiptKey"`
	IssueNumber            string       `json:"issueNumber"`
	ReceiptUrl             string       `json:"receiptUrl"`
	Amount                 float64      `json:"amount"`
	TaxFreeAmount          float64      `json:"taxFreeAmount"`
	IssueStatus            string       `json:"issueStatus"`
	Failure                *FailureInfo `json:"failure,omitempty"`
	CustomerIdentityNumber string       `json:"customerIdentityNumber"`
	RequestedAt            string       `json:"requestedAt"`
}

type DiscountInfo struct {
	Amount int `json:"amount"`
}

type RefundReceiveAccountInfo struct {
	BankCode      string `json:"bankCode"`
	AccountNumber string `json:"accountNumber"`
	HolderName    string `json:"holderName"`
}

type CashReceipt struct {
	ReceiptKey             string         `json:"receiptKey"`
	IssueNumber            string         `json:"issueNumber"`
	IssueStatus            string         `json:"issueStatus"`
	Amount                 int            `json:"amount"`
	TaxFreeAmount          int            `json:"taxFreeAmount"`
	OrderId                string         `json:"orderId"`
	OrderName              string         `json:"orderName"`
	Type                   string         `json:"type"`
	TransactionType        string         `json:"transactionType"`
	BusinessNumber         string         `json:"businessNumber"`
	CustomerIdentityNumber string         `json:"customerIdentityNumber"`
	Failure                *FailureDetail `json:"failure,omitempty"`
	RequestedAt            string         `json:"requestedAt"`
	ReceiptUrl             string         `json:"receiptUrl"`
}

type FailureDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Cancel struct {
	CancelAmount          int       `json:"cancelAmount"`
	CancelReason          string    `json:"cancelReason"`
	TaxFreeAmount         int       `json:"taxFreeAmount"`
	TaxExemptionAmount    int       `json:"taxExemptionAmount"`
	RefundableAmount      int       `json:"refundableAmount"`
	EasyPayDiscountAmount int       `json:"easyPayDiscountAmount"`
	CanceledAt            time.Time `json:"canceledAt"`
	TransactionKey        string    `json:"transactionKey"`
	ReceiptKey            string    `json:"receiptKey,omitempty"`
}
