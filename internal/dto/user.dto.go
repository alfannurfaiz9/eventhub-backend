package dto

import "time"

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

type Notification struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
