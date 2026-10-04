package dto

import (
	"mime/multipart"
	"time"
)

type User struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	ImgUrl   *string `json:"img_url"`
	Address  *string `json:"address"`
	Bio      *string `json:"bio"`
	Status   string  `json:"status"`
	Role     string  `json:"role"`
}

type UserProfile struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	ImgUrl   *string `json:"img_url"`
	Address  *string `json:"address"`
	Bio      *string `json:"bio"`
	Role     string  `json:"role"`
}

type UpdateProfile struct {
	FullName        string                `form:"full_name"`
	Img             *multipart.FileHeader `form:"img_url"`
	Address         *string               `form:"address"`
	Bio             *string               `form:"bio"`
	CurrentPassword string                `form:"current_password"`
	NewPassword     string                `form:"new_password"`
}

type UserInformation struct {
	FullName string  `json:"full_name"`
	Email    string  `json:"email"`
	ImgUrl   *string `json:"img_url"`
}

type Notification struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
