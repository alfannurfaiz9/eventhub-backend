package dto

type Speaker struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
	Company  string `json:"company"`
}
