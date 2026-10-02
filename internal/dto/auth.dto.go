package dto

type Register struct {
	FullName string `json:"full_name" example:"your name"`
	Email    string `json:"email" example:"youremail@gmail.com"`
	Password string `json:"password" example:"yourpassword"`
}

type Login struct {
	Email    string `json:"email" example:"youremail@gmail.com"`
	Password string `json:"password" example:"yourpassword"`
}

type ForgotPassword struct {
	Email       string `json:"email" example:"youremail@gmail.com"`
	NewPassword string `json:"new_password" example:"yournewpassword"`
}
