package model

type Profile struct {
	ID       int64     `json:"id,omitempty"`
	Title    string    `json:"title,omitempty"`
	Category *Category `json:"category,omitempty"`
	Deleted  bool      `json:"deleted,omitempty"`
}
