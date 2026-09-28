package dto

type Testimony struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Message  string `json:"message"`
}
