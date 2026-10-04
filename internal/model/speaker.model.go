package model

type Speaker struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	ImgUrl   string `db:"img_url"`
	Position string `db:"position"`
	Company  string `db:"company"`
}
