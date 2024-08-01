package model

type Replace struct {
	ID           int64 `json:"id,omitempty"`
	TransferFrom int64 `json:"transfer_from,omitempty"`
	TransferTo   int64 `json:"transfer_to,omitempty"`
}
