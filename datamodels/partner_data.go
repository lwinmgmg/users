package datamodels

type PartnerData struct {
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Email            string `json:"email"`
	IsEmailConfirmed bool   `json:"is_email_confirmed"`
	Phone            string `json:"phone"`
	IsPhoneConfirmed bool   `json:"is_phone_confirmed"`
	Code             string `json:"code"`
}
