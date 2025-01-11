package models

type (
	AccrualResponse struct {
		OrderNumber   string        `json:"order"`
		AccrualStatus AccrualStatus `json:"status"`
		Accrual       Money         `json:"accrual"`
	}

	AccrualStatus string
)

const (
	AccrualStatusRegistered AccrualStatus = "REGISTERED"
	AccrualStatusProcessing AccrualStatus = "PROCESSING"
	AccrualStatusInvalid    AccrualStatus = "INVALID"
	AccrualStatusProcessed  AccrualStatus = "PROCESSED"
)
