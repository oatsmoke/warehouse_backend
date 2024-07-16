package model

type Report struct {
	Categories     []*Category                     `json:"categories"`
	Departments    []*Department                   `json:"departments"`
	Leftover       map[int64][]*Location           `json:"leftover"`
	Total          map[int64][]*Location           `json:"total"`
	FromStorage    map[int64][]*Location           `json:"from_storage"`
	ToStorage      map[int64][]*Location           `json:"to_storage"`
	FromContract   map[int64][]*Location           `json:"from_contract"`
	ToContract     map[int64][]*Location           `json:"to_contract"`
	FromDepartment map[int64]map[int64][]*Location `json:"from_department"`
	ToDepartment   map[int64]map[int64][]*Location `json:"to_department"`
}
