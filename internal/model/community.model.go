package model

import "time"

type Community struct {
	Id          int       `db:"id"`
	Name        string    `db:"name"`
	ImgUrl      string    `db:"img_url"`
	Description string    `db:"description"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CommunityList struct {
	Community
	Category
	TotalMember   int `db:"total_member"`
	UpcomingEvent int `db:"upcoming_event"`
}

type CommunityMember struct {
	Id       int     `db:"id"`
	FullName string  `db:"full_name"`
	Job      *string `db:"job"`
}

type UserCommunity struct {
	UserId      int       `db:"user_id"`
	CommunityId int       `db:"community_id"`
	CreatedAt   time.Time `db:"created_at"`
}
