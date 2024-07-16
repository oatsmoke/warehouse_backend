package model

type Contract struct {
	ID        int64  `json:"id"`
	Number    string `json:"number"`
	Address   string `json:"address"`
	IsDeleted bool   `json:"is_deleted"`
}
