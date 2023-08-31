package transactions

type Transaction struct {
	SenderId   int   `json:"sender_id"`
	ReceiverId int   `json:"receiver_id"`
	Amount     int64 `json:"amount"`
}
