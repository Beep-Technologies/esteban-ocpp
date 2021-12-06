package schemas

// StartTransactionResponse
type StartTransactionResponse struct {
	IdTagInfo     *IdTagInfo `json:"idTagInfo"`
	TransactionId int        `json:"transactionId"`
}
