package webHandler

type RegisterParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfoParam struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type FeedParam struct {
	Token      string `json:"token"`
	LatestTime string `json:"latest_time"`
}
