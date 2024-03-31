package model

type Report struct {
	Categories     []Category                 `json:"categories,omitempty"`
	Departments    []Department               `json:"departments,omitempty"`
	Leftover       map[int][]Location         `json:"leftover,omitempty"`
	Total          map[int][]Location         `json:"total,omitempty"`
	FromStorage    map[int][]Location         `json:"fromStorage,omitempty"`
	ToStorage      map[int][]Location         `json:"toStorage,omitempty"`
	FromContract   map[int][]Location         `json:"fromContract,omitempty"`
	ToContract     map[int][]Location         `json:"toContract,omitempty"`
	FromDepartment map[int]map[int][]Location `json:"fromDepartment,omitempty"`
	ToDepartment   map[int]map[int][]Location `json:"toDepartment,omitempty"`
}
