package model

import (
	"time"
)

type Event struct {
	Id          int       `db:"id"`
	Title       string    `db:"title"`
	ImgUrl      *string   `db:"img_url"`
	Description string    `db:"description"`
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	Format      string    `db:"format"`
	Capacity    int       `db:"capacity"`
	OrganizerId int       `db:"organizer_id"`
	CommunityId *int      `db:"community_id"`
	LocationId  int       `db:"location_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type UserEvent struct {
	UserId    int       `db:"user_id"`
	EventId   int       `db:"event_id"`
	CreatedAt time.Time `db:"created_at"`
}

type EventCategory struct {
	EventId    int `db:"event_id"`
	CategoryId int `db:"category_id"`
}

type EventSpeaker struct {
	EventId   int `db:"event_id"`
	SpeakerId int `db:"speaker_id"`
}

type EventDiscussion struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventList struct {
	Event
	UserEvent
	Category
	EventCategory
	Community
	Location
	User
	TotalAttendee int `db:"total_attendee"`
}

type EventDetail struct {
	Event
	Category
	Location
	User
	Community
	TotalAttendee int `db:"total_attendee"`
}
