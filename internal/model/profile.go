package model

type Profile struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Category *Category `json:"category"`
	Deleted  bool      `json:"deleted"`
}
