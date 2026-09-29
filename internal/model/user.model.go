package model

import "time"

type User struct {
	Id        int       `db:"id"`
	FullName  string    `db:"full_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	ImgUrl    *string   `db:"img_url"`
	Address   *string   `db:"address"`
	Bio       *string   `db:"bio"`
	Status    string    `db:"status"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Notification struct {
	Id          int       `db:"id"`
	UserId      int       `db:"user_id"`
	ReadAt      time.Time `db:"read_at"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Type        string    `db:"type"`
	CreatedAt   time.Time `db:"created_at"`
}
