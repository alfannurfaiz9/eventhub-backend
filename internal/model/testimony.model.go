package model

import "time"

type Testimony struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	Message   string    `db:"message"`
	Company   string    `db:"company"`
	Position  string    `db:"position"`
	CreatedAt time.Time `db:"created_at"`
}

type TestimonyList struct {
	User
	Testimony
}
