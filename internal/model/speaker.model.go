package model

type Speaker struct {
	Id       int    `db:"id"`
	Name     string `db:"name"`
	Position string `db:"position"`
	Company  string `db:"company"`
}
