package dto

import "time"

type OrganizerDashboard struct {
	TotalEvent    int    `json:"total_event"`
	TotalAttendee int    `json:"total_attendee"`
	AvgFillRate   string `json:"avg_fill_rate"`
}

type AdminDashboard struct {
	TotalUser      int    `json:"total_user"`
	TotalEvent     int    `json:"total_event"`
	TotalCommunity int    `json:"total_community"`
	AvgFillRate    string `json:"avg_fill_rate"`
}

type UserList struct {
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"joined_at"`
}
