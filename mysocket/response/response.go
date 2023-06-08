package response

type UserInfo struct {
	UID      string `json:"uid" validate:"required"`
	UserName string `json:"username" validate:"required"`
	ID       string `json:"id" validate:"required"`
}
