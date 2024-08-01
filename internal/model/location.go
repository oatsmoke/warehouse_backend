package model

type Location struct {
	ID             int64       `json:"id,omitempty"`
	Date           int64       `json:"date,omitempty"`
	Code           string      `json:"code,omitempty"`
	Equipment      *Equipment  `json:"equipment,omitempty"`
	Employee       *Employee   `json:"employee,omitempty"`
	Company        *Company    `json:"company,omitempty"`
	FromDepartment *Department `json:"from_department,omitempty"`
	FromEmployee   *Employee   `json:"from_employee,omitempty"`
	FromContract   *Contract   `json:"from_contract,omitempty"`
	ToDepartment   *Department `json:"to_department,omitempty"`
	ToEmployee     *Employee   `json:"to_employee,omitempty"`
	ToContract     *Contract   `json:"to_contract,omitempty"`
	TransferType   string      `json:"transfer_type,omitempty"`
	Price          int         `json:"price,omitempty"`
}

type RequestLocation struct {
	Date         int64  `json:"date"`
	EquipmentId  int64  `json:"equipment_id"`
	Way          string `json:"way"`
	ThisLocation string `json:"this_location"`
	Where        string `json:"where"`
	InDepartment bool   `json:"in_department"`
	Company      int64  `json:"company"`
	ToDepartment int64  `json:"to_department"`
	ToEmployee   int64  `json:"to_employee"`
	ToContract   int64  `json:"to_contract"`
	TransferType string `json:"transfer_type"`
	Price        int    `json:"price"`
}

type LocationAndRequestLocation struct {
	Location        *Location          `json:"location"`
	RequestLocation []*RequestLocation `json:"request_location"`
}
