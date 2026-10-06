package dto

type Community struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ImgUrl      string `json:"img_url"`
	Description string `json:"description"`
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
	FullName string `json:"full_name"`
	Job      string `json:"job"`
}

type UserCommunity struct {
	UserId      int `json:"user_id"`
	CommunityId int `json:"community_id"`
}
