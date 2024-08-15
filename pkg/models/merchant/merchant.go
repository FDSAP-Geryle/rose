package models

import "time"

type Recordable interface {
	TableName() string
}

type Record struct {
	Id          int
	Date        time.Time
	FilePath    string
	FilePathErr string
	Status      string
	User        string
	Notes       string
}

func (Record) TableName() string {
	return "upload_merchant"
}

type Recordwip struct {
	Id          int
	Date        time.Time
	FilePath    string
	FilePathErr string
	Status      string
	User        string
	Notes       string
}

func (Recordwip) TableName() string {
	return "upload_merchant_wip"
}

type Recordok struct {
	Id           int
	Date         time.Time
	FilePath     string
	FilePathErr  string
	Status       string
	User         string
	Notes        string
	TotalUpload  int
	TotalSuccess int
	TotalError   int
}

func (Recordok) TableName() string {
	return "upload_merchant_ok"
}

type ReceiveUploadActivated struct {
	MerchantName  string    `json:"Merchant Name"`
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
	AgentName     string    `json:"Agent Name"`
	AgentId       int       `json:"Agent Id"`
}

func (ReceiveUploadActivated) TableName() string {
	return "receive_upload_activated"
}
