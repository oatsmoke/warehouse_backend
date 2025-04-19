package model

import "time"

type Location struct {
	ID             int64       `json:"id,omitempty"`
	Date           *time.Time  `json:"date,omitempty"`
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
	Price          string      `json:"price,omitempty"`
}

type RequestLocation struct {
	Date         *time.Time `json:"date,omitempty"`
	EquipmentId  int64      `json:"equipment_id,omitempty"`
	Way          string     `json:"way,omitempty"`
	ThisLocation string     `json:"this_location,omitempty"`
	Where        string     `json:"where,omitempty"`
	InDepartment bool       `json:"in_department,omitempty"`
	Company      int64      `json:"company,omitempty"`
	ToDepartment int64      `json:"to_department,omitempty"`
	ToEmployee   int64      `json:"to_employee,omitempty"`
	ToContract   int64      `json:"to_contract,omitempty"`
	TransferType string     `json:"transfer_type,omitempty"`
	Price        string     `json:"price,omitempty"`
}

type LocationAndRequestLocation struct {
	Location        *Location          `json:"location"`
	RequestLocation []*RequestLocation `json:"request_location"`
}
