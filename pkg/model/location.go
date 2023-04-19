package model

type Location struct {
	Id           int        `json:"id,omitempty" db:"id"`
	Date         int64      `json:"date,omitempty" db:"date"`
	Code         string     `json:"code,omitempty" db:"code"`
	Equipment    Equipment  `json:"equipment" db:"equipment"`
	Employee     Employee   `json:"employee" db:"employee"`
	Company      Company    `json:"company" db:"company"`
	ToDepartment Department `json:"toDepartment" db:"to_department"`
	ToEmployee   Employee   `json:"toEmployee" db:"to_employee"`
	ToContract   Contract   `json:"toContract" db:"to_contract"`
	TransferType string     `json:"transferType,omitempty" db:"transfer_type"`
	Price        string     `json:"price,omitempty" db:"price"`
}

type RequestLocation struct {
	EquipmentId  int    `json:"equipmentId,omitempty"`
	ThisLocation string `json:"thisLocation,omitempty"`
	Date         int64  `json:"date,omitempty"`
	Where        string `json:"where,omitempty"`
	InDepartment bool   `json:"inDepartment,omitempty"`
	Company      int    `json:"company,omitempty"`
	ToDepartment int    `json:"toDepartment,omitempty"`
	ToEmployee   int    `json:"toEmployee,omitempty"`
	ToContract   int    `json:"toContract,omitempty"`
	TransferType string `json:"transferType,omitempty"`
	Price        string `json:"price,omitempty"`
}

type LocationAndRequestLocation struct {
	Location        Location          `json:"location"`
	RequestLocation []RequestLocation `json:"requestLocation"`
}
