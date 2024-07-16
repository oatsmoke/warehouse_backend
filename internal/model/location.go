package model

type Location struct {
	ID             int64       `json:"id"`
	Date           int64       `json:"date"`
	Code           string      `json:"code"`
	Equipment      *Equipment  `json:"equipment"`
	Employee       *Employee   `json:"employee"`
	Company        *Company    `json:"company"`
	FromDepartment *Department `json:"from_department"`
	FromEmployee   *Employee   `json:"from_employee"`
	FromContract   *Contract   `json:"from_contract"`
	ToDepartment   *Department `json:"to_department"`
	ToEmployee     *Employee   `json:"to_employee"`
	ToContract     *Contract   `json:"to_contract"`
	TransferType   string      `json:"transfer_type"`
	Price          int         `json:"price"`
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
