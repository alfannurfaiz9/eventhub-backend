package model

type OrganizerDashboard struct {
	TotalEvent    int    `db:"total_event"`
	TotalAttendee int    `db:"total_attendee"`
	AvgFillRate   string `db:"avg_fill_rate"`
}

type AdminDashboard struct {
	TotalUser      int    `db:"total_user"`
	TotalEvent     int    `db:"total_event"`
	TotalCommunity int    `db:"total_community"`
	AvgFillRate    string `db:"avg_fill_rate"`
}
