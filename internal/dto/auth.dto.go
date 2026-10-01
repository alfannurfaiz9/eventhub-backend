package dto

type Register struct {
	FullName string `json:"full_name" example:"your name"`
	Email    string `json:"email" example:"youemail@gmail.com"`
	Password string `json:"password" example:"yourpassword"`
}

type Login struct {
	Email    string `json:"email" example:"youemail@gmail.com"`
	Password string `json:"password" example:"yourpassword"`
}
