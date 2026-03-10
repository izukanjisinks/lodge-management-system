package models

// ProfileResponse combines user and employee data for profile display
type ProfileResponse struct {
	User     *User     `json:"user"`
	Employee *Employee `json:"employee,omitempty"`
}
