package response

import models "rose/pkg/models/struct"

type ResponseModel struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type EncryptDecyrptResponseModel struct {
	RetCode string                         `json:"retCode"`
	Message string                         `json:"message"`
	Data    []models.EncryptDecryptRequest `json:"data"`
}
