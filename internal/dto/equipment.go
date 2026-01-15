package dto

type Equipment struct {
	CompanyID    int64  `json:"company_id,omitempty" binding:"required"`
	ProfileID    int64  `json:"profile_id,omitempty" binding:"required"`
	SerialNumber string `json:"serial_number,omitempty" binding:"required"`
}

type CreateEquipmentRequest struct {
	Date          string   `json:"date,omitempty" binding:"required"`
	CompanyID     int64    `json:"company_id,omitempty" binding:"required"`
	ProfileID     int64    `json:"profile_id,omitempty" binding:"required"`
	SerialNumbers []string `json:"serial_numbers,omitempty" binding:"required"`
	Param         string   `json:"param,omitempty" binding:"required"`
	ParamID       int64    `json:"param_id,omitempty"`
}
