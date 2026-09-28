package dto

import "time"

type Event struct {
	Title    string    `json:"title"`
	ImgUrl   string    `json:"img_url"`
	Category string    `json:"category"`
	StartAt  time.Time `json:"start_at"`
}

type EventList struct {
	Title         string    `json:"title"`
	ImgUrl        string    `json:"img_url"`
	Category      string    `json:"category"`
	StartAt       time.Time `json:"start_at"`
	Location      string    `json:"location"`
	TotalAttendee int       `json:"total_attendee"`
	Capacity      int       `json:"capacity"`
}

type EventDetail struct {
	Title         string    `json:"title"`
	ImgUrl        string    `json:"img_url"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	StartAt       time.Time `json:"start_at"`
	Location      string    `json:"location"`
	TotalAttendee int       `json:"total_attendee"`
	Capacity      int       `json:"capacity"`
	Organizer     string    `json:"organizer"`
	Community     string    `json:"community"`
}

type UserEvent struct {
	UserId  int `json:"user_id"`
	EventId int `json:"event_id"`
}
