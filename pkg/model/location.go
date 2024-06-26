package model

type Location struct {
	Id             int        `json:"id,omitempty" db:"id"`
	Date           int64      `json:"date,omitempty" db:"date"`
	Code           string     `json:"code,omitempty" db:"code"`
	Equipment      Equipment  `json:"equipment" db:"equipment"`
	Employee       Employee   `json:"employee" db:"employee"`
	Company        Company    `json:"company" db:"company"`
	ToDepartment   Department `json:"toDepartment" db:"to_department"`
	ToEmployee     Employee   `json:"toEmployee" db:"to_employee"`
	ToContract     Contract   `json:"toContract" db:"to_contract"`
	FromDepartment Department `json:"fromDepartment" db:"from_department"`
	FromEmployee   Employee   `json:"fromEmployee" db:"from_employee"`
	FromContract   Contract   `json:"fromContract" db:"from_contract"`
	TransferType   string     `json:"transferType,omitempty" db:"transfer_type"`
	Price          int        `json:"price,omitempty" db:"price"`
}

type RequestLocation struct {
	Date         int64  `json:"date,omitempty"`
	EquipmentId  int    `json:"equipmentId,omitempty"`
	Way          string `json:"way,omitempty"`
	ThisLocation string `json:"thisLocation,omitempty"`
	Where        string `json:"where,omitempty"`
	InDepartment bool   `json:"inDepartment,omitempty"`
	Company      int    `json:"company,omitempty"`
	ToDepartment int    `json:"toDepartment,omitempty"`
	ToEmployee   int    `json:"toEmployee,omitempty"`
	ToContract   int    `json:"toContract,omitempty"`
	TransferType string `json:"transferType,omitempty"`
	Price        int    `json:"price,omitempty"`
}

type LocationAndRequestLocation struct {
	Location        Location          `json:"location"`
	RequestLocation []RequestLocation `json:"requestLocation"`
}
