package dto

import "time"

type Community struct {
	Id          int     `json:"id"`
	Name        *string `json:"name"`
	ImgUrl      string  `json:"img_url"`
	Description string  `json:"description"`
}

type CommunityList struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"img_url"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	TotalMember   int    `json:"total_member"`
	UpcomingEvent int    `json:"upcoming_event"`
}

type CommunityMember struct {
	Id       int     `json:"id"`
	FullName string  `json:"full_name"`
	ImgUrl   *string `json:"img_url"`
}

type UserCommunity struct {
	UserId      int `json:"user_id"`
	CommunityId int `json:"community_id"`
}

type CommunityDiscussion struct {
	FullName  string    `json:"full_name"`
	ImgUrl    *string   `json:"img_url"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
