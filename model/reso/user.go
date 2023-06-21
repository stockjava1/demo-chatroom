package reso

// GetUser GET "/user" response object
type GetUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Gender   int64  `json:"gender"`
	Age      int64  `json:"age"`
	Interest string `json:"interest"`
}

// PutUser PUT "/user" response object
type PutUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Gender   int64  `json:"gender"`
	Age      int64  `json:"age"`
	Interest string `json:"interest"`
}

// PostUser POST "/user" response object
type PostUser struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

// PostLogin POST "/login" response object
type PostLogin struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	Token    string `json:"token"`
}
