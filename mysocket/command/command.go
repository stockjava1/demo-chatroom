package command

type Login struct {
	token string `json:"token" validate:"required,gte=10"`
}

type Logout struct {
}

type User struct {
}
