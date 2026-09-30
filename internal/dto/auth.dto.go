package dto

type Register struct {
	FullName string `json:"full_name" example:"marianus example"`
	Email    string `json:"email" example:"example@gmail.com"`
	Password string `json:"password" example:"password"`
}

type Login struct {
	Email    string `json:"email" example:"example@gmail.com"`
	Password string `json:"password" example:"examplepass"`
}
