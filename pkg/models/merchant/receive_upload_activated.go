package merchant

import "time"

type ReceiveUploadActivated struct {
	PhoneNumber   string    `json:"Phone Number"`
	OwnerName     string    `json:"Owner Name"`
	IDType        string    `json:"ID Type"`
	IDNumber      string    `json:"ID Number"`
	DOB           time.Time `json:"DoB"`
	BIN           string    `json:"BIN"`
	Cid           string    `json:"CID"`
	AccountNumber string    `json:"Account Number"`
	AccountType   string    `json:"Account Type"`
	AccountName   string    `json:"Account Name"`
	UploadedAt    time.Time `json:"created_at"`
	OwnerAddress  string    `json:"Owner Address"`
	BranchName    string    `json:"Branch Name"`
	BranchCode    string    `json:"Branch Code"`
	UnitName      string    `json:"Unit Name"`
	UnitCode      string    `json:"Unit Code"`
	CenterName    string    `json:"Center Name"`
	CenterCode    string    `json:"Center Code"`
	Mid           string    `json:"MID"`
	Mpan          string    `json:"MPAN"`
}

func (ReceiveUploadActivated) TableName() string {
	return "receive_upload_activated"
}

type TempMerchantOk struct {
	BIN        string `json:"BIN"`
	Cid        string `json:"CID"`
	BranchName string `json:"Branch Name"`
	BranchCode string `json:"Branch Code"`
	UnitName   string `json:"Unit Name"`
	UnitCode   string `json:"Unit Code"`
	CenterName string `json:"Center Name"`
	CenterCode string `json:"Center Code"`
	AOCode     string `json:"AO Code"`
	AOName     string `json:"AO Name"`
}

func (TempMerchantOk) TableName() string {
	return "temp_merchant_ok"
}

type ListUploadMerchant struct {
	No        string `json:"No"`
	StoreName string `json:"Store Name"`
}
