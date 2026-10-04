package dto

import "time"

type Event struct {
	Title       string    `json:"title"`
	ImgUrl      *string   `json:"img_url"`
	Description string    `json:"description"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Format      string    `json:"format"`
	Capacity    int       `json:"capacity"`
	OrganizerId int       `json:"organizer_id"`
	CommunityId *int      `json:"community_id"`
	LocationId  int       `json:"location_id"`
}

type CreateEvent struct {
	Title           string    `form:"title"`
	ImgUrl          *string   `form:"img_url"`
	Description     string    `form:"description"`
	StartAt         time.Time `form:"start_at"`
	EndAt           time.Time `form:"end_at"`
	Format          string    `form:"format"`
	Capacity        int       `form:"capacity"`
	CommunityId     *int      `form:"community_id"`
	CategoryId      int       `form:"category_id"`
	LocationName    string    `form:"location_name"`
	SpeakerName     string    `form:"speaker_name"`
	SpeakerImgUrl   string    `form:"speaker_img_url"`
	SpeakerPosition string    `form:"speaker_position"`
	SpeakerCompany  string    `form:"speaker_company"`
}

type EventCategory struct {
	CategoryId int `json:"category_id"`
}

type EventSpeaker struct {
	EventId   int `db:"event_id"`
	SpeakerId int `db:"speaker_id"`
}

type EventList struct {
	Title         string    `json:"title"`
	ImgUrl        *string   `json:"img_url"`
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

type JoinEvent struct {
	EventId int `json:"event_id"`
}
