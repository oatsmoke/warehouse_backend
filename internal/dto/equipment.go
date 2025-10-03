package dto

type Equipment struct {
	SerialNumber string `json:"serial_number,omitempty" binding:"required"`
	ProfileID    int64  `json:"profile_id,omitempty" binding:"required"`
}
