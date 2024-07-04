package models

type EncryptDecryptRequest struct {
	Host     string `json:"host"`
	DbName   string `json:"dbName"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type EncryptDecryptRequestV2 struct {
	Text string `json:"text"`
}
