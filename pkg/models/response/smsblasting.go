package response

import models "rose/pkg/models/struct"

type SMSContentResponse struct {
	RetCode string              `json:"retCode"`
	Message string              `json:"message"`
	Data    []models.SMSContent `json:"data"`
}

type DistinctSMSContentResponse struct {
	RetCode string      `json:"retCode"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type EMPCDistinctSMSTypeResponse struct {
	RetCode string                   `json:"retCode"`
	Message string                   `json:"message"`
	Data    []models.DistinctSMSType `json:"data"`
}
