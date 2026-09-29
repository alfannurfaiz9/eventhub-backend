package dto

type OrganizerDashboard struct {
	TotalEvent    int    `json:"total_event"`
	TotalAttendee int    `json:"total_attendee"`
	AvgFillRate   string `json:"avg_fill_rate"`
}
