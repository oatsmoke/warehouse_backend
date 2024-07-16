package model

type Replace struct {
	ID           int64 `json:"id"`
	TransferFrom int64 `json:"transfer_from"`
	TransferTo   int64 `json:"transfer_to"`
}
