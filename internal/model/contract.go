package model

type Contract struct {
	ID      int64  `json:"id"`
	Number  string `json:"number"`
	Address string `json:"address"`
	Deleted bool   `json:"deleted"`
}
