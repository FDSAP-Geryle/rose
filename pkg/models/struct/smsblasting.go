package models

type SMSContent struct {
	SmsId      int    `json:"smsId"`
	SmsType    string `json:"smsType"`
	SmsContent string `json:"smsContent"`
}

type DistinctSMSType struct {
	SmsType string `json:"smsType"`
}
