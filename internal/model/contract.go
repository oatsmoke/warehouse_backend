package model

type Contract struct {
	ID      int64  `json:"id,omitempty"`
	Number  string `json:"number,omitempty"`
	Address string `json:"address,omitempty"`
	Deleted bool   `json:"deleted,omitempty"`
}
